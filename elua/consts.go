package elua

import (
	"github.com/gogf/gf/v2/container/gmap"
	"reflect"
	"strings"
)

var (
	localTypeFunc = gmap.NewStrAnyMap(true)
	//localStructFunc = gmap.NewStrAnyMap[string, *any](true)

	//localFunc = gmap.NewStrAnyMap(true)

	localFuncString = gmap.NewStrStrMap(true)
)

func RegisterFuncString(name string, script string) {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	localFuncString.Set(name, script)
}

func RegisterTypeFunc(name string, value any) {
	if reflect.ValueOf(value).Kind() != reflect.Struct {
		panic("value must be struct")
		return
	}
	localTypeFunc.Set(name, value)
}

//func RegisterStructFunc(name string, value any) {
//	if reflect.ValueOf(value).Kind() != reflect.Struct {
//		panic("value must be struct")
//		return
//	}
//	localStructFunc.Set(name, &value)
//}

//func RegisterFunc(name string, value any) {
//	if reflect.ValueOf(value).Kind() != reflect.Func {
//		panic("value must be struct")
//		return
//	}
//	localFunc.Set(name, &value)
//}
