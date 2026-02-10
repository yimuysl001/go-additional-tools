package pkgs

import (
	"context"
	"github.com/Masterminds/sprig/v3"
	"github.com/gogf/gf/v2/container/gmap"
	timeconv "github.com/Andrew-M-C/go.timeconv"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
	"time"
)

var localData = gmap.NewStrAnyMap(true)

func DefPkgs() {
	common()
	sprigFunc()
}

func common() {
	RegisterImport("common", map[string]any{
		"addDate": timeconv.AddDate,
		"getLocalData": func() *gmap.StrAnyMap {
			return localData
		},
		"setCommon": func(key string, value any) {
			localData.Set(key, value)
		},
		"setCommonByFunc": func(key string, f func() interface{}) {
			localData.GetOrSetFuncLock(key, f)
		},
		"getCommon": func(key string) (any, bool) {
			return localData.Search(key)
		},

		"setCommonCache": func(ctx context.Context, key string, value string, d time.Duration) error {
			return gcache.Set(ctx, key, value, d)
		},
		"getCommonCache": func(ctx context.Context, key string) (*gvar.Var, error) {
			return gcache.Get(ctx, key)
		},
	})
}

func sprigFunc() {
	RegisterImport("sprig", sprig.FuncMap())
}
