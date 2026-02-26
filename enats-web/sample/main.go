package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"go-additional-tools/econf"
	"go-additional-tools/enats-web/enats"
	"go-additional-tools/enats-web/eserver"
	"time"
)

func main() {
	econf.MustInitConf() // 初始化配置

	//go enats.NatsServer() // 启用nats
	time.Sleep(time.Second)
	natsWebSample() // 启用测试模块

	httpweb() // 统一网关
}

func httpweb() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		eserver.GfRouter(group)
	})
	s.Run()
}

func natsWebSample() {

	s := eserver.NewServer("testNats")
	s.Use(eserver.AddMiddleware("test/bbb", func(ctx *enats.MsgContext, next enats.HandlerFunc) {
		g.Log().Info(ctx, "test middleware bbb")
		next(ctx)
	}), eserver.AddMiddleware("test", func(ctx *enats.MsgContext, next enats.HandlerFunc) {
		g.Log().Info(ctx, "test middleware1")
		next(ctx)
	}),
		eserver.AddMiddleware("test", func(ctx *enats.MsgContext, next enats.HandlerFunc) {
			next(ctx)
			g.Log().Info(ctx, "test middleware2")
		}),
		eserver.AddMiddleware("test", func(ctx *enats.MsgContext, next enats.HandlerFunc) {
			g.Log().Info(ctx, "test middleware3")
			next(ctx)
			g.Log().Info(ctx, "test middleware3")
		}),
	)

	s.Get("/test/*act", func(ctx *enats.MsgContext) error {
		g.Log().Info(ctx, "test request")
		g.Log().Info(ctx, "GetMsgMethod", ctx.GetMsgMethod())
		g.Log().Info(ctx, "GetMsgUrlPath", ctx.GetMsgUrlPath())
		g.Log().Info(ctx, "Params", ctx.GetMsgParam())
		ctx.SetResponseBody([]byte("接收成功"))
		return nil
	})

	s.Get("/static/pdf", func(ctx *enats.MsgContext) error {

		ctx.SetResponseBody(gfile.GetBytes("E:\\data\\xwechat_files\\wxid_rj8zox9d5ol222_5f9f\\msg\\file\\2026-02\\Coati 介绍.pdf"))

		return nil
	})

	s.Run()
	g.Log().Info(gctx.GetInitCtx(), "start auth server")
}
