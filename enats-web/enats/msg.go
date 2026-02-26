package enats

import (
	"context"
	"encoding/base64"
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
		subj    = r.GetRouter("subj").String()
		method  = r.Method
		url     = r.GetRouter("act").String()
		headers = r.Request.Header
		param   = r.GetRequestMap()
		body    = r.GetBody()
	)

	msg := nats.NewMsg(webPre + subj)
	for k, v := range headers {
		for i, v := range v {
			if i == 0 {
				SetHeaderString(msg, k, v)
			} else {
				AddHeaderString(msg, k, v)
			}
		}
	}

	// 添加请求方法到消息头
	SetHeaderString(msg, RequestMethod, method)

	// 添加请求路径到消息头
	SetHeaderString(msg, RequestUrlPath, url)

	// 添加请求参数到消息头
	SetHeaderString(msg, RequestParam, gjson.New(param).MustToJsonString())

	// 添加请求追踪ID到消息头
	value := ctx.Value(RespondTraceId)
	if value != nil {
		SetHeaderString(msg, RespondTraceHeader, gconv.String(value))
	} else {
		SetHeaderString(msg, RespondTraceHeader, gctx.CtxId(ctx))
	}
	SetHeaderString(msg, RespondConsumerKey, g.Cfg().MustGet(ctx, consumerKeyPath, "").String())

	msg.Data = body
	return msg
}

func SetHeaderString(msg *nats.Msg, key string, valus string) {
	msg.Header.Set(key, base64.StdEncoding.EncodeToString([]byte(valus)))
}
func AddHeaderString(msg *nats.Msg, key string, valus string) {
	msg.Header.Add(key, base64.StdEncoding.EncodeToString([]byte(valus)))
}
func GetHeaderString(msg *nats.Msg, key string) string {
	get := msg.Header.Get(key)
	if get == "" {
		return ""
	}
	decodeString, err := base64.StdEncoding.DecodeString(get)
	if err != nil {
		return get
	}
	return string(decodeString)
}
func GetMsgParam(msg *nats.Msg) *gjson.Json {
	return gjson.New(GetHeaderString(msg, RequestParam))
}

func GetMsgData(msg *nats.Msg) []byte {
	return msg.Data
}

func GetMsgMethod(msg *nats.Msg) string {
	return GetHeaderString(msg, RequestMethod)
}

func GetMsgUrlPath(msg *nats.Msg) string {
	return GetHeaderString(msg, RequestUrlPath)
}

func GetMsgTraceId(msg *nats.Msg) string {
	return GetHeaderString(msg, RespondTraceHeader)
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
	m := make(map[string][]string)
	for k, _ := range msg.Header {
		m[k] = []string{GetMsgHeaderString(msg, k)}
	}
	return m
}
func GetMsgHeaderString(msg *nats.Msg, key string) string {
	return GetHeaderString(msg, key)
}

// NewMsg 创建带跟踪ID的消息
func NewMsg(ctx context.Context, subj string, data []byte) *nats.Msg {
	msg := nats.NewMsg(subj)
	msg.Data = data
	// 添加请求追踪ID到消息头
	value := ctx.Value(RespondTraceId)
	if value != nil {
		SetHeaderString(msg, RespondTraceHeader, gconv.String(value))
	} else {
		SetHeaderString(msg, RespondTraceHeader, gctx.CtxId(ctx))
	}
	SetHeaderString(msg, RespondConsumerKey, g.Cfg().MustGet(ctx, consumerKeyPath, "").String())
	return msg

}

func SetResponseHeader(msg *nats.Msg, key, value string) {
	SetHeaderString(msg, key, value)

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
	return GetHeaderString(msg, RespondErrorHeader)
}

func GetStatusMsg(msg *nats.Msg) string {
	status := GetHeaderString(msg, RespondStatus)
	if status == "" {
		status = "200"
	}
	return status
}
