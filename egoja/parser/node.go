package eparser

import (
	"github.com/gogf/gf/v2/util/gconv"
	"go-additional-tools/egoja"

	"strings"
)

var (
	CONCAT_TOKEN_PARSER   = NewGenericTokenParser("${", "}", false)
	REPLACE_TOKEN_PARSER  = NewGenericTokenParser("#{", "}", true)
	IF_TOKEN_PARSER       = NewGenericTokenParser("?{", "}", true)
	IF_PARAM_TOKEN_PARSER = NewGenericTokenParser("?{", ",", true)
)

const defaultPlaceholder = "?"

func ParseSql(sqlstr string, varMap map[string]any, placeholder ...string) (sql string, parameters []any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	var ph = defaultPlaceholder
	if len(placeholder) > 0 && placeholder[0] != "" {
		ph = placeholder[0]
	}

	sql = strings.TrimSpace(sqlstr)
	parameters = make([]any, 0)
	sql = IF_TOKEN_PARSER.Parse(sql, func(text string) string {
		var ifTrue = false
		val := IF_PARAM_TOKEN_PARSER.Parse("?{"+text, func(s string) string {
			simple, err := egoja.ExecSimple(s, varMap)
			if err != nil {
				panic(gconv.String(err) + ":" + s)
			}
			ifTrue = simple.ToBoolean() //   gconv.Bool(simple)
			return ""
		})
		if ifTrue {
			return val
		}
		return ""
	})

	sql = CONCAT_TOKEN_PARSER.Parse(sql, func(text string) string {
		simple, err := egoja.ExecSimple(text, varMap)
		if err != nil {
			panic(gconv.String(err) + ":" + text)
		}
		return simple.String()
	})

	sql = REPLACE_TOKEN_PARSER.Parse(sql, func(text string) string {
		simple, err := egoja.ExecSimple(text, varMap)
		if err != nil {
			panic(gconv.String(err) + ":" + text)
		}
		if simple == nil || simple.String() == "null" || simple.String() == "undefined" {
			parameters = append(parameters, nil)
			return ph
		}
		interfaces := gconv.Interfaces(simple.Export())
		if interfaces == nil || len(interfaces) == 0 {
			parameters = append(parameters, simple.Export())
			return ph
		}
		parameters = append(parameters, interfaces...)
		var strs = make([]string, len(interfaces))
		for i, _ := range interfaces {
			strs[i] = ph
		}
		return strings.Join(strs, ",")
	})

	return sql, parameters, nil
}
