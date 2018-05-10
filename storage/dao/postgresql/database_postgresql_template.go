package postgresql

import (
	"database/sql"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/lib/pq"
)

func (db *DatabasePostgreSQL) GetAllTemplates() ([]*model.Template, error) {
	q := `
		SELECT u.id, u.code, u.created_at, u.updated_at
		FROM schema.template u
	`
	rows, err := db.session.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	us := make([]*model.Template, 0)
	for rows.Next() {
		u := model.Template{}
		err := rows.Scan(&u.ID, &u.Code, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}

func (db *DatabasePostgreSQL) GetTemplatesByID(id string) (*model.Template, error) {
	q := `
		SELECT u.id, u.code, u.created_at, u.updated_at
		FROM schema.template u
		WHERE u.id = $1
	`
	row := db.session.QueryRow(q, id)

	u := model.Template{}
	err := row.Scan(&u.ID, &u.Code, &u.CreatedAt, &u.UpdatedAt)
	if errPq, ok := err.(*pq.Error); ok {
		return nil, handlePgError(errPq)
	}
	if err == sql.ErrNoRows {
		return nil, dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	return &u, err
}

func (db *DatabasePostgreSQL) CreateTemplate(template *model.Template) error {
	q := `
		INSERT INTO schema.template
			(code)
		VALUES
			($1)
		RETURNING id, created_at
	`

	err := db.session.
		QueryRow(q, template.Code).
		Scan(&template.ID, &template.CreatedAt)
	if errPq, ok := err.(*pq.Error); ok {
		return handlePgError(errPq)
	}
	return err
}

func (db *DatabasePostgreSQL) DeleteTemplate(id string) error {
	q := `
		DELETE FROM schema.template
		WHERE id = $1
	`

	_, err := db.session.Exec(q, id)
	if errPq, ok := err.(*pq.Error); ok {
		return handlePgError(errPq)
	}
	return err
}

func (db *DatabasePostgreSQL) UpdateTemplate(template *model.Template) error {
	q := `
		UPDATE schema.template
		SET
			code = $2
		WHERE id = $1
		RETURNING updated_at
	`

	err := db.session.
		QueryRow(q, template.ID, template.Code).
		Scan(&template.UpdatedAt)
	if errPq, ok := err.(*pq.Error); ok {
		return handlePgError(errPq)
	}
	return err
}
