package forward

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"go-additional-tools/econf"
	"go-additional-tools/enats-web/enats"
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

func ForwardWeb(opts ...config.Option) {
	econf.MustInitConf() // 初始化配置

	natsw := enats.Nats("http")

	var optsc = []config.Option{
		server.WithHostPorts(g.Cfg().MustGet(gctx.GetInitCtx(), "server.address").String()),
		server.WithTracer(&ServerTracer{}),
	}

	if len(opts) > 0 {
		optsc = append(optsc, opts...)
	}

	h := server.Default(optsc...)
	newClient, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	h.Any("/:subj/*action", func(ctx context.Context, c *app.RequestContext) {

		var subj = c.Param("subj")
		var action = c.Param("action")
		g.Log().Info(ctx, "subj:{}", subj)
		g.Log().Info(ctx, "action:{}", action)
		msg, err2 := natsw.Request(ctx, webNatsSubPre+subj, []byte(action), 5*time.Second)
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
