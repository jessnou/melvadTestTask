package mocks

import "database/sql"

type MockPostgres struct {
	ExecFunc func(query string, args ...interface{}) (sql.Result, error)
	GetFunc  func(dest interface{}, query string, args ...interface{}) error
}

func (p MockPostgres) Exec(query string, args ...interface{}) (sql.Result, error) {
	if p.ExecFunc != nil {
		return p.ExecFunc(query, args...)
	}
	return nil, nil
}

func (p MockPostgres) Get(dest interface{}, query string, args ...interface{}) error {
	if p.GetFunc != nil {
		return p.GetFunc(dest, query, args...)
	}
	return nil
}
