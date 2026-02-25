package enats

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nats-io/nats.go"
)

// NewMsg 创建带跟踪ID的消息
func NewMsgByRequest(r *ghttp.Request) *nats.Msg {
	var (
		ctx     = r.Context()
		subj    = r.Get("subj").String()
		method  = r.Method
		url     = r.Get("act").String()
		headers = r.Request.Header
		param   = r.GetRequestMap()
		body    = r.GetBody()
	)

	msg := nats.NewMsg(webPre + subj)
	for k, v := range headers {
		for i, v := range v {
			if i == 0 {
				msg.Header.Set(k, v)
			} else {
				msg.Header.Add(k, v)
			}
		}
	}

	// 添加请求方法到消息头
	msg.Header.Set(RequestMethod, method)

	// 添加请求路径到消息头
	msg.Header.Set(RequestUrlPath, url)

	// 添加请求参数到消息头
	msg.Header.Set(RequestParam, gjson.New(param).MustToJsonString())

	// 添加请求追踪ID到消息头
	value := ctx.Value(RespondTraceId)
	if value != nil {
		msg.Header.Set(RespondTraceHeader, gconv.String(value))
	} else {
		msg.Header.Set(RespondTraceHeader, gctx.CtxId(ctx))
	}
	msg.Header.Set(RespondConsumerKey, g.Cfg().MustGet(ctx, consumerKeyPath, "").String())

	msg.Data = body
	return msg
}

func GetMsgParam(msg *nats.Msg) *gjson.Json {
	return gjson.New(msg.Header.Get(RequestParam))
}

func GetMsgData(msg *nats.Msg) []byte {
	return msg.Data
}

func GetMsgMethod(msg *nats.Msg) string {
	return msg.Header.Get(RequestMethod)
}

func GetMsgUrlPath(msg *nats.Msg) string {
	return msg.Header.Get(RequestUrlPath)
}

func GetMsgTraceId(msg *nats.Msg) string {
	return msg.Header.Get(RespondTraceHeader)
}
func GetMsgCtx(msg *nats.Msg) context.Context {
	ctx := context.WithValue(gctx.New(), RespondTraceId, GetMsgTraceId(msg))
	ctx = context.WithValue(ctx, RequestMsg, GetMsgData(msg))
	return ctx
}
func GetRequestMsg(ctx context.Context) *nats.Msg {
	msg := ctx.Value(RequestMsg)
	if msg == nil {
		return nil
	}
	return msg.(*nats.Msg)

}

func GetMsgHeader(msg *nats.Msg) map[string][]string {
	return msg.Header
}
func GetMsgHeaderString(msg *nats.Msg, key string) string {
	return msg.Header.Get(key)
}

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
	msg.Header.Set(RespondConsumerKey, g.Cfg().MustGet(ctx, consumerKeyPath, "").String())
	return msg

}

func SetResponseHeader(msg *nats.Msg, key, value string) {
	msg.Header.Set(key, value)

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

func GetStatusMsg(msg *nats.Msg) string {
	status := msg.Header.Get(RespondStatus)
	if status == "" {
		status = "200"
	}
	return status
}
