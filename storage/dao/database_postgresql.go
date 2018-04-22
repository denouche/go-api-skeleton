package dao

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	pgCodeUniqueViolation     = "23505"
	pgCodeForeingKeyViolation = "23503"
)

func handlePgError(e *pq.Error) error {
	if e.Code == pgCodeUniqueViolation {
		return newDAOError(ErrTypeDuplicate, e)
	}

	if e.Code == pgCodeForeingKeyViolation {
		return newDAOError(ErrTypeForeignKeyViolation, e)
	}
	return e
}

type DatabasePostgreSQL struct {
	session *sql.DB
}

func NewDatabasePostgreSQL(connectionURI string) Database {
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
