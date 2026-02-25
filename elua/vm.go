package elua

import (
	"context"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"strings"
	"time"
)

func ExecScript(ctx context.Context, script string, params map[string]any, d ...time.Duration) (v any, err error) {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("ctx", luar.New(L, ctx))
	var canle context.CancelFunc
	if len(d) > 0 && d[0] > 0 {
		ctx, canle = context.WithTimeout(ctx, d[0])
		defer canle()
		L.SetContext(ctx)
	}

	if params != nil && len(params) > 0 {
		L.SetGlobal("params", luar.New(L, params))
	}
	var back = make(map[string]any, 1)
	L.SetGlobal("back", luar.New(L, back))
	//localFunc.Iterator(func(k string, v any) bool {
	//	L.SetGlobal(k, luar.New(L, v))
	//	return true
	//})

	localTypeFunc.Iterator(func(k string, v any) bool {
		L.SetGlobal(k, luar.NewType(L, v))
		return true
	})
	err = L.DoString(flashScript(script))
	return back["back"], err

}

func flashScript(script string) string {

	var sb strings.Builder
	for _, s := range strings.Split(script, "\n") {
		news := strings.TrimSpace(s)
		if strings.HasPrefix(news, "return ") {
			s = strings.Replace(s, "return ", "back.back = ", 1) + "\n return back"

		}
		sb.WriteString(s)

		sb.WriteString("\n")

	}

	return sb.String()
}
