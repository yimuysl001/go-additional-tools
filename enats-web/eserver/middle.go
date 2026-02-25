package eserver

import (
	"go-additional-tools/enats-web/enats"
	"strings"
)

type Middleware func(ctx *enats.MsgContext, next enats.HandlerFunc)

type MiddlewareStruct struct {
	Name   string
	Func   Middleware
	Except []string
}

func AddMiddleware(name string, middleware Middleware, except ...string) MiddlewareStruct {
	name = strings.ReplaceAll(name, "\\", "/")
	name = strings.Trim(name, "/")
	name = name + "/"
	exceptName := make([]string, len(except))
	for i, v := range except {
		exceptName[i] = strings.ReplaceAll(v, "\\", "/")
		exceptName[i] = strings.Trim(exceptName[i], "/")
		exceptName[i] = "/" + exceptName[i]
	}

	return MiddlewareStruct{
		Name:   name,
		Func:   middleware,
		Except: exceptName,
	}
}
