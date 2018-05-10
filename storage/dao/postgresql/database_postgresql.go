package postgresql

import (
	"database/sql"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	pgCodeUniqueViolation     = "23505"
	pgCodeForeingKeyViolation = "23503"
)

func handlePgError(e *pq.Error) error {
	if e.Code == pgCodeUniqueViolation {
		return dao.NewDAOError(dao.ErrTypeDuplicate, e)
	}

	if e.Code == pgCodeForeingKeyViolation {
		return dao.NewDAOError(dao.ErrTypeForeignKeyViolation, e)
	}
	return e
}

type DatabasePostgreSQL struct {
	session *sql.DB
}

func NewDatabasePostgreSQL(connectionURI string) dao.Database {
	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to get a connection to the postgres db")
	}
	err = db.Ping()
	if err != nil {
		logrus.WithError(err).Fatal("Unable to ping the postgres db")
	}
	return &DatabasePostgreSQL{session: db}
}
