package eswagger

import (
	"context"
	"encoding/base64"
	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gconv"
	"path"
	"strings"
	"time"
)

func InitWeb(ctx context.Context, name ...any) bool {
	var n = ""
	if len(name) > 0 && name[0] != "" {
		n = gconv.String(name[0]) + "."
	}

	serverRoot := g.Cfg().MustGet(ctx, "server."+n+"serverRoot").String()
	if serverRoot == "" { // 未配置路径，不需要处理swagger
		return false
	}
	if !gfile.Exists(serverRoot) {
		gfile.Mkdir(serverRoot)
	}

	binContent, err := gaes.Decrypt(databin, CryptoKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}
	if !strings.Contains(serverRoot, ":") { // 全路径不处理
		var pwd = gfile.Pwd()
		pwd = strings.ReplaceAll(pwd, "\\", "/")
		serverRoot = path.Join(pwd, serverRoot)
	}
	if !strings.HasSuffix(serverRoot, "/") {
		serverRoot = serverRoot + "/"
	}

	if err := gres.Add(string(binContent), serverRoot); err != nil {
		g.Log().Error(ctx, err)
		return false
	}

	g.Log().Debug(ctx, "swagger 添加完成")
	return true

}

func EnhanceOpenAPIDoc(ctx context.Context, s *ghttp.Server) *SwaggerInfo {
	openapi := s.GetOpenApi()
	//openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	//openapi.Config.CommonResponseDataField = `Data`
	var swaggerInfo = new(SwaggerInfo)

	err := g.Cfg().MustGet(ctx, "swagger."+s.GetName()).Scan(swaggerInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	// API description.
	openapi.Info = goai.Info{
		Title:          swaggerInfo.Title,
		Description:    swaggerInfo.Description,
		Version:        swaggerInfo.Version,
		TermsOfService: swaggerInfo.TermsOfService,
		Contact: &goai.Contact{
			Name:  swaggerInfo.Name,
			URL:   swaggerInfo.Url,
			Email: swaggerInfo.Email,
		},
	}
	return swaggerInfo
}
func InitSwagger(ctx context.Context, f func(ctx context.Context) *goai.OpenApiV3, name ...any) {
	var n = ""
	if len(name) > 0 && name[0] != "" {
		n = gconv.String(name[0]) + "."
	}
	openapiPath := g.Cfg().MustGet(ctx, "server."+n+"openapiPath").String()
	if openapiPath == "" {
		return
	}
	if !strings.HasPrefix(openapiPath, "/") {
		openapiPath = "/" + openapiPath
	}
	if !InitWeb(ctx, name...) {
		return
	}

	s := g.Server(name...)
	var info = EnhanceOpenAPIDoc(ctx, s)

	//s.BindHookHandler("/doc/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
	//	g.Log().Debug(r.Context(), "test")
	//	flag := authSwagger(r, info)
	//	if !flag {
	//		r.ExitAll()
	//	}
	//	r.Middleware.Next()
	//
	//})

	if f != nil {
		s.BindHandler(apijson, func(r *ghttp.Request) {
			authSwagger(r, info)
			r.Response.Header().Set("connection", "keep-alive")
			r.Response.Header().Set("content-type", "application/json;charset=UTF-8")

			//r.Response.RedirectTo(g.Cfg().MustGet(r.Context(), "server.openapiPath").String())
			r.Response.Write(f(r.Context()))
		})
	}

	if openapiPath != apijson {
		s.BindHandler(apijson, func(r *ghttp.Request) {
			authSwagger(r, info)
			r.Response.RedirectTo(openapiPath)
		})
	}

}

func authSwagger(r *ghttp.Request, s *SwaggerInfo) (flag bool) {

	if len(s.Auths) <= 0 {
		return true
	}

	defer func() {
		if !flag {
			r.Response.Status = 401
			r.Response.Header().Set("WWW-Authenticate", `Basic realm="Swagger API Documentation"`)
			r.ExitAll()
		}
	}()
	cook := r.Cookie.Get("swagger-auth")

	if !cook.IsEmpty() && cook.String() == generateSessionToken()+"OK" {
		return true
	}

	if cook.IsEmpty() {
		r.Cookie.SetCookie("swagger-auth", generateSessionToken(), "", "/", 30*time.Minute)
		return false
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false
	}
	authParts := strings.SplitN(auth, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Basic" {
		return false
	}

	payload, _ := base64.StdEncoding.DecodeString(authParts[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return false
	}
	username := pair[0]
	password := pair[1]

	for _, i2 := range s.Auths {

		if i2.User == username && i2.Pass+time.Now().Format("2006010215") == password {
			r.Cookie.SetCookie("swagger-auth", generateSessionToken()+"OK", "", "/", 30*time.Minute)
			return true
		}

	}

	return false
}

func generateSessionToken() string {
	return base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006010215")))
}
