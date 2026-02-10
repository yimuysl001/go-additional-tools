package egoja

import (
	"context"
	"github.com/dop251/goja"
	"go-additional-tools/egoja/require"
	"time"
)

// ExecSimple 简易数据处理
func ExecSimple(script string, params map[string]any) (goja.Value, error) {

	vm := require.GetNewVm()
	defer require.PutGoja(vm)

	for k, v := range params {
		vm.Set(k, v)
	}

	return vm.RunString(script)
}

func ExecScript(ctx context.Context, script string, params map[string]any, d ...time.Duration) (v any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	script = TransformScript(script)
	vm := require.GetNewVm()
	defer require.PutGoja(vm)
	//vm := require.NewVm()
	for k, v := range params {
		vm.Set(k, v)
	}
	vm.Set("ctx", ctx)
	if len(d) > 0 && d[0] > 0 {
		time.AfterFunc(d[0], func() {
			vm.Interrupt("time out")
		})
	}
	runString, err := vm.RunString(script)
	if err != nil {
		return nil, err
	}
	return runString.Export(), err
}

func ExecScriptById(ctx context.Context, id, script string, params map[string]any, d ...time.Duration) (v any, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	prog, err := GetCacheProgram(id, script)
	if err != nil {
		return nil, err
	}

	vm := require.GetNewVm()
	defer require.PutGoja(vm)

	for k, v := range params {
		vm.Set(k, v)
	}
	vm.Set("ctx", ctx)

	if len(d) > 0 && d[0] > 0 {
		time.AfterFunc(d[0], func() {
			vm.Interrupt("time out")
		})
	}
	result, err := vm.RunProgram(prog)
	if err != nil {
		return nil, err
	}

	return result.Export(), err
}
