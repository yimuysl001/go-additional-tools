package enats

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
)

type MsgContext struct {
	context.Context
	msg  *nats.Msg
	resp *nats.Msg
	exit bool
}

func NewMsgContext(msg *nats.Msg) *MsgContext {
	var ctxMsg = &MsgContext{
		msg: msg,
	}
	ctx := context.WithValue(gctx.New(), RespondTraceId, GetMsgTraceId(msg))
	//ctx = context.WithValue(ctx, RequestMsg, ctxMsg)
	ctxMsg.Context = ctx
	ctxMsg.resp = NewMsg(ctx, msg.Reply, nil)
	return ctxMsg
}

func (c *MsgContext) GetMsgParam() *gjson.Json {
	return gjson.New(c.msg.Header.Get(RequestParam))
}

func (c *MsgContext) GetMsgData() []byte {
	return c.msg.Data
}

func (c *MsgContext) GetMsgMethod() string {
	return c.msg.Header.Get(RequestMethod)
}

func (c *MsgContext) GetMsgUrlPath() string {
	return c.msg.Header.Get(RequestUrlPath)
}

func (c *MsgContext) GetMsgTraceId() string {
	return c.msg.Header.Get(RespondTraceHeader)
}
func (c *MsgContext) GetMsgCtx() context.Context {

	return c.Context
}

func (c *MsgContext) SetMsgCtx(ctx context.Context) {

	c.Context = ctx
}
func (c *MsgContext) GetMsgHeader() map[string][]string {
	return c.msg.Header
}

func (c *MsgContext) GetMsgHeaderString(key string) string {
	return c.msg.Header.Get(key)
}

func (c *MsgContext) SetResponseBody(data []byte) {
	c.resp.Data = data
}

func (c *MsgContext) SetResponseHeader(key string, value string) {
	c.resp.Header.Set(key, value)
}

func (c *MsgContext) SetResponseError(err error, status ...string) {
	c.resp.Data = []byte(err.Error())
	c.resp.Header.Set(RespondErrorHeader, err.Error())
	if len(status) > 0 && status[0] != "" {
		c.resp.Header.Set(RespondStatus, status[0])
	} else {
		c.resp.Header.Set(RespondStatus, "500")
	}
}

func (c *MsgContext) ResponseWrite(body []byte, status ...string) {
	c.resp.Data = body
	if len(status) > 0 && status[0] != "" {
		c.resp.Header.Set(RespondStatus, status[0])
	} else {
		c.resp.Header.Set(RespondStatus, "200")
	}
}

func (c *MsgContext) SetResponseStatus(status string) {
	c.resp.Header.Set(RespondStatus, status)
}

func (c *MsgContext) GetResponseMsg() *nats.Msg {
	return c.resp
}

func (c *MsgContext) Exit() {
	c.exit = true
}

func (c *MsgContext) IsExit() bool {
	return c.exit
}
