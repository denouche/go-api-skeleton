package postgresql

import (
	"database/sql"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/lib/pq"
)

func (db *DatabasePostgreSQL) GetAllUsers() ([]*model.User, error) {
	q := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.created_at, u.updated_at
		FROM users.user u
	`
	rows, err := db.session.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	us := make([]*model.User, 0)
	for rows.Next() {
		u := model.User{}
		err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}

func (db *DatabasePostgreSQL) GetUsersByID(id string) (*model.User, error) {
	q := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.created_at, u.updated_at
		FROM users.user u
		WHERE u.id = $1
	`
	row := db.session.QueryRow(q, id)

	u := model.User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt)
	if errPq, ok := err.(*pq.Error); ok {
		return nil, handlePgError(errPq)
	}
	if err == sql.ErrNoRows {
		return nil, dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	return &u, err
}

func (db *DatabasePostgreSQL) CreateUser(user *model.User) error {
	q := `
		INSERT INTO users.user
			(email, first_name, last_name)
		VALUES
			($1, $2, $3)
		RETURNING id, created_at
	`

	err := db.session.
		QueryRow(q, user.Email, user.FirstName, user.LastName).
		Scan(&user.ID, &user.CreatedAt)
	if errPq, ok := err.(*pq.Error); ok {
		return handlePgError(errPq)
	}
	return err
}

func (db *DatabasePostgreSQL) DeleteUser(id string) error {
	q := `
		DELETE FROM users.user
		WHERE id = $1
	`

	_, err := db.session.Exec(q, id)
	if errPq, ok := err.(*pq.Error); ok {
		return handlePgError(errPq)
	}
	return err
}

func (db *DatabasePostgreSQL) UpdateUser(user *model.User) error {
	q := `
		UPDATE users.user
		SET
			email = $2,
			first_name = $3,
			last_name = $4
		WHERE id = $1
		RETURNING updated_at
	`

	err := db.session.
		QueryRow(q, user.ID, user.Email, user.FirstName, user.LastName).
		Scan(&user.UpdatedAt)
	if errPq, ok := err.(*pq.Error); ok {
		return handlePgError(errPq)
	}
	return err
}
