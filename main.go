package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/nats-io/nats.go"
	"go-additional-tools/econf"
	"go-additional-tools/enats-web/enats"
)

// //go:generate goversioninfo -icon=icon.ico
func main() {
	econf.MustInitConf()

	//listen.StartNatsAndCron()
	server := g.Server()

	var port = "8127"

	server.SetAddr(port)

	natsx := enats.Nats("http")

	natsx.Subscribe("httpWebCloud.test", "webrequest", func(subj, queue string, msg *nats.Msg) {

		msg.Respond([]byte("http://127.0.0.1:8127/" + string(msg.Data)))
	})
	g.Log().SetCtxKeys("Request-TranceId")
	server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			get := r.Request.Header.Get("Request-TranceId")
			if get != "" {
				r.SetCtxVar("Request-TranceId", get)
			}

			r.Middleware.Next()
		})

		group.ALL("/*", func(r *ghttp.Request) {
			g.Log().Info(r.Context(), "query:", r.GetMap())
			g.Log().Info(r.Context(), "body:", r.GetBodyString())
			r.Response.Write(r.GetUrl())
		})
	})

	server.Run()

}
