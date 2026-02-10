package egoja

import (
	"fmt"
	"github.com/dop251/goja"
	"strings"
)

func JsonConversion(refresh, s string) (string, string) {
	s = strings.TrimSpace(s)
	if (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) ||
		(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) {
		s = "return  JSON.stringify(" + s + ")"
	}
	return refresh, s
}

func TransformScript(script string) string {
	var transformed []string
	var exports []string
	lines := strings.Split(script, "\n")

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

		} else {
			transformed = append(transformed, line)
		}

	}
	exports = append(exports, transformed...)
	return fmt.Sprintf(mainFuncScript, strings.Join(exports, "\n"))
}

func GetCacheProgram(id, script string) (*goja.Program, error) {
	pr, err2, _ := single.Do(id, func() (interface{}, error) {
		value, found := cacheProgram.Load(id)
		if found && value != nil {
			return value, nil
		}
		return FlushCache(id, script)
	})
	if err2 != nil {
		return nil, err2
	}
	return pr.(*goja.Program), nil
}

func FlushCache(id string, script string) (*goja.Program, error) {
	if script == "" {
		cacheProgram.Delete(id)
		return nil, nil
	}

	script = TransformScript(script)

	prog, err := goja.Compile("script:"+id, script, true)
	if err != nil {
		return nil, err
	}
	cacheProgram.Store(id, prog)
	return prog, nil
}
