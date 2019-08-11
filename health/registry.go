package health

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type Registry interface {
	DB() *sql.DB
	Logger() *logrus.Logger
}
