package mid

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
)

var (
	// DefaultPort is the default port for the connect server.
	localNatsClient = gmap.NewStrAnyMap(true)
	natsConfigs     = make(map[string]NatsConfig)
)

const (
	natsDefKey         = "default"
	natsPath           = "nats"
	RespondErrorHeader = "NATS-ERROR"
	RespondHeader      = "NATS-HEADER"
	RespondConsumerKey = "NATS-Consumer"
	RespondTraceId     = "NATS-TRACE-ID"
	RespondTraceHeader = "NATS-TRACE-ID"
)

// NewMsg 创建带跟踪ID的消息
func NewMsg(ctx context.Context, subj string, data []byte) *nats.Msg {
	msg := nats.NewMsg(subj)
	msg.Data = data
	// 添加请求追踪ID到消息头
	value := ctx.Value(RespondTraceId)
	if value != nil {
		msg.Header.Set(RespondTraceHeader, gconv.String(value))
	} else {
		msg.Header.Set(RespondTraceHeader, gctx.CtxId(ctx))
	}
	g.Log()
	return msg

}

func SetMsgCtx(msg *nats.Msg) context.Context {
	ctx := context.WithValue(gctx.New(), RespondHeader, msg.Header)

	if s := msg.Header.Get(RespondTraceHeader); s != "" {
		return context.WithValue(ctx, RespondTraceId, s)
	}

	return ctx
}

func GetProjName(msg *nats.Msg) string {
	if msg == nil {
		return ""
	}
	if msg.Header == nil {
		return ""
	}
	return msg.Header.Get(RespondConsumerKey)
}

func GetErrorMsg(msg *nats.Msg) string {
	return msg.Header.Get(RespondErrorHeader)
}
