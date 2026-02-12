package mid

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
	"strings"
	"sync"
	"time"
)

// NatsCli NATS客户端结构体
type NatsCli struct {
	ns     *nats.Conn                    // NATS连接对象
	jsName string                        // JetStream名称
	js     nats.JetStreamContext         // JetStream上下文
	err    error                         // 错误信息
	subMap map[string]*nats.Subscription // 订阅映射表
	mu     sync.RWMutex
}

// Nats 获取NATS客户端实例，支持多实例
func Nats(name ...string) *NatsCli {
	n := natsDefKey
	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}
	// 使用本地缓存获取或创建NATS客户端实例
	v := localNatsClient.GetOrSetFuncLock(n, func() interface{} {

		config, ok := natsConfigs[n]
		if !ok {
			// 从配置文件加载NATS配置
			get, err := g.Cfg().Get(gctx.GetInitCtx(), natsPath)
			if err != nil {
				g.Log().Error(gctx.GetInitCtx(), "get nats config error:"+err.Error())
				return &NatsCli{err: err}
			}
			err = get.Scan(&natsConfigs)
			if err != nil {
				g.Log().Error(gctx.GetInitCtx(), "read nats config error:"+err.Error())
				return &NatsCli{err: err}
			}
		}
		// 再次检查配置是否存在
		config, ok = natsConfigs[n]
		if !ok {
			g.Log().Debug(gctx.GetInitCtx(), "nats config is null")
			return &NatsCli{err: errors.New("nats config is null")}
		}
		//g.Log().Info(gctx.GetInitCtx(), "配置数据：", config)
		var (
			opts = make([]nats.Option, 0)
			url  = strings.TrimSpace(config.Url)
		)
		url = strings.ToLower(url)
		// 格式化NATS URL
		if !strings.HasPrefix(url, "nats://") {
			url = "nats://" + strings.ReplaceAll(url, " ", "")
			url = strings.ReplaceAll(url, ",", ",nats://")
		}
		// 添加连接名称配置
		if config.Name != "" {
			opts = append(opts, nats.Name(config.Name))
		}
		// 添加Token认证
		if config.Token != "" {
			opts = append(opts, nats.Token(config.Token))
		}
		// 添加用户名密码认证
		if config.Username != "" && config.Password != "" {
			opts = append(opts, nats.UserInfo(config.Username, config.Password))
		}
		// 设置连接超时时间
		if config.Timeout > 0 {
			opts = append(opts, nats.Timeout(time.Duration(config.Timeout)*time.Second))
		}
		// 设置重连参数
		if config.ReconnectWait > 0 {
			opts = append(opts, nats.ReconnectWait(time.Duration(config.ReconnectWait)*time.Second))
			opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
				g.Log().Infof(gctx.GetInitCtx(), "NATS Reconnected to %s", nc.ConnectedUrl())
			}))
		}
		// 设置最大重连次数
		if config.MaxReconnects > 0 {
			opts = append(opts, nats.MaxReconnects(config.MaxReconnects))
		}
		// 设置连接关闭回调
		opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
			g.Log().Error(gctx.GetInitCtx(), "NATS Connection Closed")
		}))
		nc, err := nats.Connect(url, opts...)
		if err != nil {
			g.Log().Error(gctx.GetInitCtx(), "NATS Connect Error:"+err.Error())
			return &NatsCli{err: err}
		}

		// 如果配置了流名称，则初始化JetStream
		if config.StreamName != "" {
			js, err := nc.JetStream()
			if err != nil {
				g.Log().Error(gctx.GetInitCtx(), "NATS JetStream Error:"+err.Error())
				return &NatsCli{err: err}
			}
			// 创建或验证流配置
			err = addAddStream(js, config.StreamConfig, config.StreamName)
			if err != nil {
				g.Log().Error(gctx.GetInitCtx(), "NATS AddStream Error:"+err.Error())
				return &NatsCli{err: err}
			}
			return &NatsCli{ns: nc, js: js, jsName: config.StreamName, subMap: make(map[string]*nats.Subscription)}
		}
		// 返回普通NATS客户端实例
		return &NatsCli{ns: nc, subMap: make(map[string]*nats.Subscription)}

	})

	return v.(*NatsCli)
}

// addAddStream 创建或验证NATS JetStream流 todo 待完善
func addAddStream(js nats.JetStreamContext, config StreamConfig, streamName string) error {
	// 检查流是否已存在
	_, err := js.StreamInfo(streamName)
	if err != nil {

		// 检查是否是因为流不存在导致的错误
		if !errors.Is(err, nats.ErrStreamNotFound) {
			return fmt.Errorf("failed to get stream info: %w", err)
		}

		// 流不存在，创建新的持久化流
		streamConfig := &nats.StreamConfig{
			Name:                 streamName,
			Description:          config.Description,
			Subjects:             config.Subjects,
			Retention:            nats.RetentionPolicy(config.Retention),
			MaxConsumers:         config.MaxConsumers,
			MaxMsgs:              config.MaxMsgs,
			MaxBytes:             config.MaxBytes,
			Discard:              nats.DiscardPolicy(config.Discard),
			DiscardNewPerSubject: config.DiscardNewPerSubject,
			MaxAge:               config.MaxAge * time.Hour, // 设置最大消息有效期
			MaxMsgsPerSubject:    config.MaxMsgsPerSubject,
			MaxMsgSize:           config.MaxMsgSize,
			Storage:              nats.StorageType(config.Storage),
			Replicas:             config.Replicas,
			NoAck:                config.NoAck,
			Sealed:               config.Sealed,
			DenyDelete:           config.DenyDelete,
			DenyPurge:            config.DenyPurge,
			AllowRollup:          config.AllowRollup,
			Compression:          nats.StoreCompression(config.Compression),
			FirstSeq:             config.FirstSeq,
			AllowDirect:          config.AllowDirect,
			MirrorDirect:         config.MirrorDirect,
			Metadata:             config.Metadata,
			Template:             config.Template,
		}

		// 创建新流
		_, err = js.AddStream(streamConfig)
		if err != nil {
			return fmt.Errorf("failed to add stream: %w", err)
		}
	}

	return nil
}

// Request 发送请求消息并等待响应
func (n *NatsCli) Request(ctx context.Context, subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	if n.err != nil {
		return nil, n.err
	}
	msg := NewMsg(ctx, subj, data)
	return n.ns.RequestMsg(msg, timeout)
}

func (n *NatsCli) PublishRequest(ctx context.Context, subj string, data []byte, count int, f func(ctx context.Context, msg *nats.Msg) error) error {
	if n.err != nil {
		return n.err
	}
	if count == 0 {
		return errors.New("无消费者")
	}
	if count == 1 {
		request, err := n.Request(ctx, subj, data, time.Second)
		if err != nil {
			return err
		}
		return f(ctx, request)
	}

	inbox := nats.NewInbox()
	sub, err := n.ns.SubscribeSync(inbox)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	msg := NewMsg(ctx, subj, data)
	msg.Reply = inbox
	err = n.ns.PublishMsg(msg)
	if err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		nextMsg, err := sub.NextMsg(time.Second)
		if err != nil {
			return err
		}
		err = f(ctx, nextMsg)
		if err != nil {
			return err
		}

	}
	return err
}

// SubscribeRequest 订阅请求主题并处理响应
func (n *NatsCli) SubscribeRequest(subj, queue string, f func(subj, queue string, msg *nats.Msg) (context.Context, []byte, error)) error {
	if n.err != nil {
		return n.err
	}
	// 检查是否已经订阅过此主题+队列组合
	_, ok := n.CheckSubscribe(subj, queue)
	if ok {
		return errors.New("订阅已存在：" + subj + ":" + queue)
	}
	// 创建队列订阅
	subscribe, err := n.ns.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		// 执行用户定义的处理函数
		ctx, bytes, err := f(subj, queue, msg)
		natsmsg := NewMsg(ctx, msg.Reply, bytes)
		// 如果处理函数返回错误，在响应消息头中设置错误信息
		if err != nil {
			natsmsg.Header.Set(RespondErrorHeader, err.Error())
		}

		value := ctx.Value(RespondHeader)

		if value != nil {
			switch m := value.(type) {
			case string:
				natsmsg.Header.Set(RespondHeader, m)
			case []byte:
				natsmsg.Header.Set(RespondHeader, string(m))
			case map[string]string:
				for k, v := range m {
					natsmsg.Header.Set(k, v)
				}
			case nats.Header:
				for k, v := range m {
					for i, vv := range v {
						if i > 0 {
							natsmsg.Header.Add(k, vv)
						} else {
							natsmsg.Header.Set(k, vv)
						}
					}
				}
			case map[string]any:
				for s, a := range m {
					natsmsg.Header.Set(s, gconv.String(a))
				}

			}

		}

		err = msg.RespondMsg(natsmsg)
		if err != nil {
			g.Log().Error(ctx, "处理出错：", err)
		}
	})
	n.setSubscribeMap(subj, queue, subscribe)
	return err

}

func (n *NatsCli) Publish(ctx context.Context, subj string, data []byte, header map[string]string, stream ...bool) error {
	if n.err != nil {
		return n.err
	}
	natsmsg := NewMsg(ctx, subj, data)
	if len(header) > 0 {
		for k, v := range header {
			natsmsg.Header.Set(k, v)
		}
	}
	if n.js != nil && len(stream) > 0 && stream[0] {
		_, err := n.js.PublishMsg(natsmsg)
		return err
	}
	return n.ns.PublishMsg(natsmsg)
}
func (n *NatsCli) Subscribe(subj, queue string, f func(subj, queue string, msg *nats.Msg), stream ...bool) error {
	if n.err != nil {
		return n.err
	}
	_, ok := n.CheckSubscribe(subj, queue)
	if ok {
		return errors.New("订阅已存在")
	}
	if n.js != nil && len(stream) > 0 && stream[0] {
		subscribe, err := n.js.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
			f(subj, queue, msg)
		})
		if err != nil {
			return err
		}
		n.setSubscribeMap(subj, queue, subscribe)
		return err
	}

	subscribe, err := n.ns.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		f(subj, queue, msg)
	})

	if err != nil {
		return err
	}
	n.setSubscribeMap(subj, queue, subscribe)
	return nil

}
func (n *NatsCli) CheckSubscribe(subj, queue string) (*nats.Subscription, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	s, ok := n.subMap[subj+":"+queue]
	return s, ok
}

func (n *NatsCli) CloseSubscribe(subj, queue string) error {
	subscribe, ok := n.CheckSubscribe(subj, queue)
	if !ok {
		return errors.New("订阅不存在")
	}
	return subscribe.Drain()

}

func (n *NatsCli) setSubscribeMap(subj, queue string, subscribe *nats.Subscription) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	n.subMap[subj+":"+queue] = subscribe
}
