package dao

import (
	"github.com/denouche/go-api-skeleton/storage/model"
)

type Database interface {
	GetAllUsers() ([]*model.User, error)
	GetUsersByID(userID string) (*model.User, error)
	CreateUser(user *model.User) error
	DeleteUser(userID string) error
	UpdateUser(user *model.User) error
}
