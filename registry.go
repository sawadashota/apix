package apix

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	db     *sql.DB
	logger *logrus.Logger
}

func newRegistry() (*Registry, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	return &Registry{
		db:     db,
		logger: newLogger(),
	}, nil
}

func (r *Registry) DB() *sql.DB {
	return r.db
}

func (r *Registry) Logger() *logrus.Logger {
	return r.logger
}
