package eswagger

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gres"
	"go-additional-tools/econf"
	"testing"
)

func TestSwagger(t *testing.T) {
	econf.MustInitConf()
	InitSwagger(gctx.GetInitCtx(), nil)
	server := g.Server()

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.WriteJson(g.Map{"code": 200, "msg": "success"})
		})
	})

	gres.Add()

	server.Run()

}
