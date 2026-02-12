package estatic

import (
	"embed"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"go-additional-tools/econf"
	"testing"
)

//go:embed doc/*
var staticFiles embed.FS

func TestFs(t *testing.T) {
	econf.MustInitConf()
	// 创建子文件系统，去除doc前缀
	//subFS, err := fs.Sub(staticFiles, "doc")
	//if err != nil {
	//	panic(err)
	//}
	InitPublic(gctx.New(), staticFiles, "doc")

	server := g.Server()

	server.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.WriteJson(g.Map{"code": 200, "msg": "success"})
		})
	})

	server.Run()

}
