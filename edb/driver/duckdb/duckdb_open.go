package duckdb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// Open creates and returns an underlying sql.DB object for duckSql.
// https://duckdb.org/
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	source, err := configNodeToSource(config)
	if err != nil {
		return nil, err
	}
	underlyingDriverName := "duckdb"
	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}

func configNodeToSource(config *gdb.ConfigNode) (string, error) {
	var source string

	if config.Link != "" {
		source = config.Link
	} else {
		source = config.Name
	}

	if config.Extra != "" {
		extraMap, err := gstr.Parse(config.Extra)
		if err != nil {
			return "", gerror.WrapCodef(
				gcode.CodeInvalidParameter,
				err,
				`invalid extra configuration: %s`, config.Extra,
			)
		}

		for k, v := range extraMap {
			if strings.Contains(source, "?") {
				source += fmt.Sprintf(`?%s=%s`, k, v)
			} else {
				source += fmt.Sprintf(`&%s=%s`, k, v)
			}

		}
	}
	return source, nil
}
