package fake

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/denouche/go-api-skeleton/utils"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/satori/go.uuid"
)

const (
	cacheKeyUsers = "users"
)

func (db *DatabaseFake) saveUsers(users []*model.User) {
	data := make([]interface{}, 0)
	for _, v := range users {
		data = append(data, v)
	}
	db.save(cacheKeyUsers, data)
}

func (db *DatabaseFake) loadUsers() []*model.User {
	users := make([]*model.User, 0)
	b, err := db.Cache.Get(cacheKeyUsers)
	if err != nil {
		return users
	}
	err = json.Unmarshal(b, &users)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Error while unmarshal fake users")
	}
	return users
}

func (db *DatabaseFake) GetAllUsers() ([]*model.User, error) {
	return db.loadUsers(), nil
}

func (db *DatabaseFake) GetUsersByID(userID string) (*model.User, error) {
	users := db.loadUsers()
	for _, u := range users {
		if u.ID == userID {
			return u, nil
		}
	}
	return nil, dao.NewDAOError(dao.ErrTypeNotFound, errors.New("user not found"))
}

func (db *DatabaseFake) CreateUser(user *model.User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()

	users := db.loadUsers()
	users = append(users, user)
	db.saveUsers(users)
	return nil
}

func (db *DatabaseFake) DeleteUser(userID string) error {
	users := db.loadUsers()
	newUsers := make([]*model.User, 0)
	for _, u := range users {
		if u.ID != userID {
			newUsers = append(newUsers, u)
		}
	}
	db.saveUsers(newUsers)
	return nil
}

func (db *DatabaseFake) UpdateUser(user *model.User) error {
	users := db.loadUsers()
	var foundUser *model.User
	for _, u := range users {
		if u.ID == user.ID {
			foundUser = u
			break
		}
	}

	if foundUser == nil {
		return dao.NewDAOError(dao.ErrTypeNotFound, errors.New("user not found"))
	}

	foundUser.UserEditable = user.UserEditable
	now := time.Now()
	foundUser.UpdatedAt = &now
	db.saveUsers(users)

	*user = *foundUser
	return nil
}
