package dao

import (
	"time"

	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/satori/go.uuid"
)

type DatabaseMock struct {
}

func NewDatabaseMock() Database {
	return &DatabaseMock{}
}

func (db *DatabaseMock) GetAllUsers() ([]*model.User, error) {
	u := make([]*model.User, 0)
	return u, nil
}

func (db *DatabaseMock) CreateUser(user *model.User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()
	return nil
}
