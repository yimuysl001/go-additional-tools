package eserver

import (
	"errors"
	"go-additional-tools/enats-web/enats"
	"strings"
)

type Server struct {
	serverName     string
	name           string
	client         *enats.NatsCli
	middlewareList []MiddlewareStruct
	methodRoute    map[string][]RoutePath
}

func NewServer(serverName string, name ...string) *Server {
	if serverName == "" {
		panic("serverName 不能为空")
	}

	n := "http"
	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}
	return &Server{
		name:           n,
		serverName:     serverName,
		client:         enats.Nats(n),
		middlewareList: make([]MiddlewareStruct, 0),
		methodRoute:    make(map[string][]RoutePath),
	}
}

func (s *Server) Get(path string, handler enats.HandlerFunc) {
	s.RegisterRoute("GET", path, handler)
}

func (s *Server) Post(path string, handler enats.HandlerFunc) {
	s.RegisterRoute("POST", path, handler)
}

func (s *Server) RegisterRoute(method string, path string, handler enats.HandlerFunc) {
	method = strings.ToUpper(method)
	path = strings.TrimSpace(path)
	path = strings.ReplaceAll(path, "\\", "")
	path = strings.Trim(path, "/")
	if path == "" {
		panic("path 不能为空")
	}
	m, ok := s.methodRoute[method]
	if !ok {
		m = []RoutePath{make(RoutePath), make(RoutePath), make(RoutePath)}
	}

	if strings.Contains(path, "*") {
		m[2][path] = handler // 通用路由最后匹配
	} else if strings.Contains(path, ":") || strings.Contains(path, "{") {
		m[1][path] = handler // 占位其次匹配
	} else {
		m[0][path] = handler // 静态路由优先匹配
	}
	s.methodRoute[method] = m
}

func (s *Server) GetRoute(method string, path string) (enats.HandlerFunc, bool) {
	method = strings.ToUpper(method)
	path = strings.TrimSpace(path)

	hs, ok := s.methodRoute[method]

	if !ok {
		return nil, false
	}

	for _, v := range hs {
		for k, f := range v {
			if ok, _ := matchRoute(k, path); ok {
				return f, true
			}
		}

	}

	if method != "ALL" {
		return s.GetRoute("ALL", path)
	}

	return nil, false
}

func (s *Server) Use(middleware ...MiddlewareStruct) {
	s.middlewareList = append(s.middlewareList, middleware...)
}
func (s *Server) buildHandler(middlewares []Middleware, final enats.HandlerFunc) enats.HandlerFunc {
	if len(middlewares) == 0 {
		return final
	}
	return func(ctx *enats.MsgContext) error {
		if ctx.IsExit() {
			return nil
		}
		var err error
		first := middlewares[0]
		rest := middlewares[1:]
		first(ctx, func(c *enats.MsgContext) error {
			if c.IsExit() {
				return nil
			}
			return s.buildHandler(rest, final)(c)
		})
		return err
	}
}

func (s *Server) handlerFunc() enats.HandlerFunc {
	return func(ctx *enats.MsgContext) error {
		f, ok := s.GetRoute(ctx.GetMsgMethod(), ctx.GetMsgUrlPath())
		if !ok {
			ctx.SetResponseError(errors.New("未找到路由"), "404")
			return nil
		}
		var applicableMiddlewares []Middleware
	w:
		for _, v := range s.middlewareList {
			for _, except := range v.Except {
				if strings.HasSuffix(ctx.GetMsgUrlPath(), except) {
					continue w
				}
			}
			if v.Name == "" {
				applicableMiddlewares = append(applicableMiddlewares, v.Func)
				continue
			}
			if !strings.HasPrefix(ctx.GetMsgUrlPath(), v.Name) {
				continue
			}
			applicableMiddlewares = append(applicableMiddlewares, v.Func)
		}

		finalHandler := func(ctx *enats.MsgContext) error {
			if ctx.IsExit() {
				return nil
			}
			return f(ctx)
		}

		handler := s.buildHandler(applicableMiddlewares, finalHandler)
		return handler(ctx)

	}
}

func (s *Server) Run() error {
	return s.client.WebSubscribe(s.serverName, s.handlerFunc())
}
