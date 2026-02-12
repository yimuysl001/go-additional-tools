package web

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"go-additional-tools/ecron/model"
	"go-additional-tools/ecron/server/listen/service"
	"go-additional-tools/ewebstatic/estatic"
)

func CronWeb() {
	//econf.MustInitConf()

	s := g.Server()
	estatic.InitPublic(gctx.New(), staticFiles, "page")
	s.Group("/api/v1/cron", func(group *ghttp.RouterGroup) {
		group.GET("/list", func(r *ghttp.Request) {
			cron, err := service.GetCronService().ListCron(r.Context())
			if err != nil {
				g.Log().Error(r.Context(), err)
				r.Response.WriteJson(g.Map{
					"code": 500,
					"msg":  "获取列表失败",
				})
				return
			}
			r.Response.WriteJson(g.Map{
				"code": 200,
				"data": cron,
			})
		})
		group.POST("/delete", func(r *ghttp.Request) {

			model := model.CronJob{
				ID: "~KEY:" + r.Get("id").String(),
			}

			_, err := service.GetCronService().DeleteCron(r.Context(), model)

			if err != nil {
				g.Log().Error(r.Context(), err)
				r.Response.WriteJson(g.Map{
					"code": 500,
					"msg":  "删除失败",
				})
				return
			}
			r.Response.WriteJson(g.Map{
				"code": 200,
				"msg":  "删除成功",
			})

		})

	})

	s.Run()

}
