package mock

import (
	"github.com/denouche/go-api-skeleton/storage/model"
)

func (db *DatabaseMock) GetAllUsers() ([]*model.User, error) {
	args := db.Called()
	return args.Get(0).([]*model.User), args.Error(1)
}

func (db *DatabaseMock) GetUsersByID(id string) (*model.User, error) {
	args := db.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

func (db *DatabaseMock) CreateUser(user *model.User) error {
	args := db.Called()
	return args.Error(0)
}

func (db *DatabaseMock) DeleteUser(id string) error {
	args := db.Called()
	return args.Error(0)
}

func (db *DatabaseMock) UpdateUser(user *model.User) error {
	args := db.Called()
	return args.Error(0)
}
