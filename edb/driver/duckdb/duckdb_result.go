package duckdb

import "database/sql"

type Result struct {
	sql.Result
	affected          int64
	lastInsertId      int64
	lastInsertIdError error
}

func (dbr Result) RowsAffected() (int64, error) {
	return dbr.affected, nil
}

func (dbr Result) LastInsertId() (int64, error) {
	return dbr.lastInsertId, dbr.lastInsertIdError
}
