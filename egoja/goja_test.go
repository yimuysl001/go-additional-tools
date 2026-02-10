package egoja

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
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
