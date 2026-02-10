package require

import (
	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/frame/g"
	"sync"
)

var (
	vmPool = sync.Pool{
		New: func() interface{} {
			//vm := goja.New()
			//vm.Set("imports", Import)
			//vm.Set("importFunc", ImportFunc)
			//vm.Set("db", dbutil.DB)

			return NewVm()
		},
	}
	localCommonParameter = make(map[string]any)
)

func NewVm() *goja.Runtime {
	vm := goja.New()
	vm.Set("require", Require)
	vm.Set("glogs", g.Log)
	vm.Set("glog", g.Log())
	for k, v := range localCommonParameter {
		vm.Set(k, v)

	}
	return vm
}

func GetNewVm() *goja.Runtime {
	return vmPool.Get().(*goja.Runtime) //  vm
}
func PutGoja(vm *goja.Runtime) {
	vm.ClearInterrupt()
	vmPool.Put(vm)
}

// 注册通用变量
func RegisterCommonParameter(name string, value any) {
	localCommonParameter[name] = value
}

// 注册通用变量
func RemoveCommonParameter(name string) {
	delete(localCommonParameter, name)
}
