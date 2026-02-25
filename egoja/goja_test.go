package egoja

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/os/gctx"
	"go-additional-tools/egoja/pkgs"
	"go-additional-tools/egoja/require"
	"sync"
	"testing"
)

var scripts = map[string]string{
	"mathUtils": `
            export function add(a, b) { return a + b; }
            export function multiply(a, b) { return a * b; }
        `,
	"dataProcessor": `
            import { add, multiply } from 'mathUtils.js';
            export function process(data) {
                return data.map(item => multiply(item.value, add(item.base, 1)));
            };
			export function processOne(data) {
                return multiply(data.value, add(data.base, 1));
            };
        `,
}

func TestScript(t *testing.T) {

	for name, source := range scripts {
		require.RegisterFuncScript(name, source)
	}
	data := make([]map[string]any, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]any{"value": i + 1, "base": i + 2}
	}
	var wg sync.WaitGroup
	for i, datum := range data {
		wg.Add(1)
		item := datum
		go func() {
			defer wg.Done()

			id, err := ExecScriptById(gctx.New(), "test", `import {processOne} from 'dataProcessor.js'
return processOne(data)`, map[string]any{"data": item})
			fmt.Println(i, "结果：", id, err)
		}()

	}
	wg.Wait()
}

func TestScript2(t *testing.T) {

	for name, source := range scripts {
		require.RegisterFuncScript(name, source)
	}
	data := make([]map[string]any, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]any{"value": i + 1, "base": i + 2}
	}
	var wg sync.WaitGroup
	for i, datum := range data {
		wg.Add(1)
		item := datum
		go func() {
			defer wg.Done()
			id, err := ExecScript(gctx.New(), `import {processOne} from 'dataProcessor.js'
return processOne(data)`, map[string]any{"data": item})
			fmt.Println(i, "结果：", id, err)
		}()

	}
	wg.Wait()
}

func TestScript3(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := ExecScriptById(gctx.New(), "test", `
return index+1
`, map[string]any{"index": i})
			fmt.Println("结果：", id, err)
		}()
	}
	wg.Wait()
}

func TestScript4(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := ExecScriptById(gctx.New(), "test", `
var ccc =  index+1
return  ctx.Value("vm").Get("ccc")
`, map[string]any{"index": i})
			fmt.Println("结果：", id, err)
		}()
	}
	wg.Wait()
}

func TestScript5(t *testing.T) {
	//	id, err := ExecScriptById(gctx.New(), "test", `
	//  var foo = "bar";
	//return  ctx.Value("vm").Get("foo")
	//`, map[string]any{"index": 123})
	//	fmt.Println("结果：", id, err)

	vm := goja.New()

	// 注册一个函数，通过闭包持有 vm
	vm.Set("getGlobal", func(call goja.FunctionCall) goja.Value {
		name := call.Arguments[0].String()
		val := vm.Get(name) // 使用闭包捕获的 vm
		return val
	})

	// 脚本中调用
	a, e := vm.RunString(`(function() {
 				  var foo = "bar";
    return getGlobal("foo"); // 输出 "bar"
        })()
  
`)
	fmt.Println(a, e)

}

func TestGoScript(t *testing.T) {
	pkgs.GoEnv()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := ExecScriptById(gctx.New(), "testgo", `
import * as strings from 'strings'
var sb1 = strings.Builder
var sb2 = strings.Builder
sb1.WriteString("123==>")

for (var i = 0; i < 10; i++) { 
  sb1.WriteString((index+i)+'')
 }

sb2.WriteString("=================\n")
sb2.WriteString(sb1.String())
return sb2.String()
`, map[string]any{"index": i})
			fmt.Println("结果：", id, err)
		}()
	}
	wg.Wait()
}

func BenchmarkVariableP(b *testing.B) {
	require.RegisterCommonParameter("targs", func(b string, args ...any) {
		fmt.Println(b)
		fmt.Println(args...)
	})

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := ExecScriptById(context.TODO(), "AAA", `var aaa= ["a","b","ccc",ddd]
targs("a",  ...aaa)
`, map[string]any{
				"ddd": "123546",
			})
			fmt.Println(err)
		}()

	}
	wg.Wait()

}
