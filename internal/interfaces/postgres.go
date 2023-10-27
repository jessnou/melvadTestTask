package interfaces

import "database/sql"

type Postgres interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
}
