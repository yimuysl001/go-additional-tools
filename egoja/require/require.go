package require

import (
	"fmt"
	"github.com/dop251/goja"
	"go-additional-tools/egoja/pkgs"
	"strings"
	"sync"
)

var (
	//cacheProgram = sync.Map{} todo 添加后异步处理会又问题
	cacheScript = sync.Map{}
)

func Require(name string) any {
	name = strings.ToLower(strings.TrimSpace(name))

	if !strings.HasSuffix(name, ".js") { // todo 先从缓存中获取 如果没有去找脚本
		cache, b := pkgs.GetCache(name)
		if b {
			return cache
		}
		//panic("require: " + name + " is not a valid module name")
		name = name + ".js"
	}

	//value, ok := cacheProgram.Load(name)
	//if ok {
	//	return value
	//}

	script, ok := cacheScript.Load(name)
	if !ok {
		panic("cacheScript require: " + name + " is not a valid module name")
	}

	compile, err := goja.Compile("func:"+name, script.(string), false)
	if err != nil {
		panic("compile require: " + name + " is not a valid module name")
	}

	vm := NewVm()
	exportsObj := make(map[string]any)
	vm.Set("exports", exportsObj)
	_, err = vm.RunProgram(compile)
	if err != nil {
		panic("runProgram require: " + name + " " + err.Error())
	}
	//exportsObj := v.Export()
	//cacheProgram.Store(name, exportsObj)
	return exportsObj

}

func RemoveScript(name string) {
	cacheScript.Delete(name)
	//cacheProgram.Delete(name)
}

func RegisterFuncScript(name string, script string) {
	name = strings.ToLower(strings.TrimSpace(name))
	if !strings.HasSuffix(name, ".js") {
		name = name + ".js"
	}
	cacheScript.Store(name, transformExports(script))
	//cacheProgram.Delete(name)
}

func transformExports(script string) string {
	lines := strings.Split(script, "\n")
	var transformed []string
	var exports []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "import ") {
			trimmed = strings.TrimPrefix(trimmed, "import ")
			parts := strings.SplitN(trimmed, " from ", 2)
			imports := strings.TrimSpace(parts[0])
			modulePath := strings.Trim(strings.TrimSpace(parts[1]), "'\";")

			if imports == "* as " || strings.HasPrefix(imports, "* as ") {
				varName := strings.TrimPrefix(imports, "* as ")
				exports = append(exports, fmt.Sprintf("const %s = require('%s');", varName, modulePath))
			} else {
				// 默认导入
				exports = append(exports, fmt.Sprintf("const %s = require('%s');", imports, modulePath))
			}
		} else if strings.HasPrefix(trimmed, "export function ") {
			funcName := strings.TrimPrefix(trimmed, "export function ")
			funcNames := strings.SplitN(strings.TrimSpace(funcName), "(", 2)
			transformed = append(transformed, "exports."+funcNames[0]+"= function ("+funcNames[1])
		} else if strings.HasPrefix(trimmed, "function ") {
			funcName := strings.TrimPrefix(trimmed, "function ")
			funcNames := strings.SplitN(strings.TrimSpace(funcName), "(", 2)
			transformed = append(transformed, "exports."+funcNames[0]+"= function ("+funcNames[1])
		} else if strings.HasPrefix(trimmed, "export const ") {
			funcName := strings.TrimPrefix(trimmed, "export const ")
			transformed = append(transformed, "exports."+strings.TrimSpace(funcName))
		} else if strings.HasPrefix(trimmed, "const ") {
			funcName := strings.TrimPrefix(trimmed, "const ")
			transformed = append(transformed, "exports."+strings.TrimSpace(funcName))
		} else {
			transformed = append(transformed, line)
		}

	}

	exports = append(exports, transformed...)
	return strings.Join(exports, "\n")
}
