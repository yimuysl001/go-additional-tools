package duckdb

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by duckdb driver`,
		)

	case gdb.InsertOptionDefault:
		tableFields, err := d.GetCore().GetDB().TableFields(ctx, table)
		if err == nil {
			for _, field := range tableFields {
				if field.Key == "pri" {
					pkField := *field
					ctx = context.WithValue(ctx, internalPrimaryKeyInCtx, pkField)
					break
				}
			}
		}
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}
