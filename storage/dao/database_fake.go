package dao

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/allegro/bigcache"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	cacheKeyUsers = "users"
)

type DatabaseFake struct {
	Cache *bigcache.BigCache
}

func (db *DatabaseFake) save(users []*model.User) {
	b, err := json.Marshal(users)
	if err != nil {
		logrus.WithError(err).Error("Error while marshal fake users")
		db.Cache.Set(cacheKeyUsers, []byte("[]"))
		return
	}
	err = db.Cache.Set(cacheKeyUsers, b)
	if err != nil {
		logrus.WithError(err).Error("Error while saving fake users")
	}
}

func (db *DatabaseFake) load() []*model.User {
	users := make([]*model.User, 0)
	b, err := db.Cache.Get(cacheKeyUsers)
	if err != nil {
		return users
	}
	err = json.Unmarshal(b, &users)
	if err != nil {
		logrus.WithError(err).Error("Error while unmarshal fake users")
	}
	return users
}

func NewDatabaseFake() Database {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	if err != nil {
		logrus.WithError(err).Fatal("Error while instantiate cache")
	}
	return &DatabaseFake{
		Cache: cache,
	}
}

func (db *DatabaseFake) GetAllUsers() ([]*model.User, error) {
	return db.load(), nil
}

func (db *DatabaseFake) GetUsersByID(userID string) (*model.User, error) {
	users := db.load()
	for _, u := range users {
		if u.ID == userID {
			return u, nil
		}
	}
	return nil, newDAOError(ErrTypeNotFound, errors.New("user not found"))
}

func (db *DatabaseFake) CreateUser(user *model.User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()

	users := db.load()
	users = append(users, user)
	db.save(users)
	return nil
}

func (db *DatabaseFake) DeleteUser(userID string) error {
	users := db.load()
	newUsers := make([]*model.User, 0)
	for _, u := range users {
		if u.ID != userID {
			newUsers = append(newUsers, u)
		}
	}
	db.save(newUsers)
	return nil
}

func (db *DatabaseFake) UpdateUser(user *model.User) error {
	users := db.load()
	var foundUser *model.User
	for _, u := range users {
		if u.ID == user.ID {
			foundUser = u
			break
		}
	}

	if foundUser == nil {
		return newDAOError(ErrTypeNotFound, errors.New("user not found"))
	}

	foundUser.UserEditable = user.UserEditable
	now := time.Now()
	foundUser.UpdatedAt = &now
	db.save(users)

	*user = *foundUser
	return nil
}
