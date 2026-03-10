package main

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/nats-io/nats.go"
	"go-additional-tools/econf"
	"go-additional-tools/enats-web/enats"
	"go-additional-tools/enats-web/eserver"
	"time"
)

type ServerTracer struct {
	// contain entities which recording metric
}

// Start record the beginning of an RPC invocation.
func (s *ServerTracer) Start(ctx context.Context, _ *app.RequestContext) context.Context {
	// do nothing
	return gctx.New()
}

// Finish record after receiving the response of server.
func (s *ServerTracer) Finish(ctx context.Context, c *app.RequestContext) {

}

func main() {
	econf.MustInitConf() // 初始化配置
	//
	////go enats.NatsServer() // 启用nats
	//time.Sleep(time.Second)
	//natsWebSample() // 启用测试模块
	//
	//httpweb() // 统一网关
	natsw := enats.Nats("http")

	h := server.Default(server.WithHostPorts(":8080"),
		server.WithTracer(&ServerTracer{}),

	)
	newClient, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	h.Any("/:subj/*action", func(ctx context.Context, c *app.RequestContext) {

		var subj = c.Param("subj")
		var action = c.Param("action")
		g.Log().Info(ctx, "subj:{}", subj)
		g.Log().Info(ctx, "action:{}", action)
		msg, err2 := natsw.Request(ctx, "httpWebCloud."+subj, []byte(action), 10*time.Second)

		if err2 != nil {
			if errors.Is(err2, nats.ErrNoResponders) {
				c.WriteString("暂无接收者")
			} else {
				g.Log().Error(ctx, "nats error", err2)
				c.WriteString(err2.Error())
			}
			return

		}

		g.Log().Info(ctx, "action:{}", action)

		g.Log().Info(ctx, "msg.Data:{}", string(msg.Data))

		s := c.QueryArgs().String()
		requestURI := string(msg.Data) //+ "/" + action
		if s != "" {
			requestURI = requestURI + "?" + s
		}

		c.Request.SetRequestURI(requestURI)
		c.Request.Header.Set("Request-TranceId", gctx.CtxId(ctx))
		err2 = newClient.Do(ctx, &c.Request, &c.Response)
		if err2 != nil {
			g.Log().Error(ctx, "http error", err2)
			c.WriteString(err2.Error())
		}

	})
	h.Run()

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
