package elua

import (
	luar "layeh.com/gopher-luar"
	"strings"
)

func Require(l *luar.LState) int {

	checkInt := l.CheckString(1)
	if checkInt == "" {
		l.RaiseError("require error: no parameters were passed in. ")
	}

	checkInt = strings.ToLower(checkInt)
	alias := ""

	n := strings.SplitN(checkInt, " as ", 2)

	if len(n) == 2 {
		alias = strings.TrimSpace(n[1])
		checkInt = strings.TrimSpace(n[0])
	} else {
		checkInt = strings.TrimSpace(checkInt)
	}

	value, ok := localFuncString.Search(checkInt)

	if !ok {
		l.RaiseError("require error 【%s】:no data was retrieved by the script.", checkInt)
	}
	err := l.DoString(value)

	if err != nil {
		l.RaiseError("require error 【%s】: %s", checkInt, err.Error())
	}
	if alias != "" {
		l.SetGlobal(alias, l.Get(-1))
		l.Pop(1)
	}

	return 1
}
