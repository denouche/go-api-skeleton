package fake

import (
	"time"

	"encoding/json"

	"github.com/allegro/bigcache"
	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/sirupsen/logrus"
)

type DatabaseFake struct {
	Cache *bigcache.BigCache
}

func NewDatabaseFake() dao.Database {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	if err != nil {
		logrus.WithError(err).Fatal("Error while instantiate cache")
	}
	return &DatabaseFake{
		Cache: cache,
	}
}

func (db *DatabaseFake) save(key string, data []interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Errorf("Error while marshal fake %s", key)
		db.Cache.Set(key, []byte("[]"))
		return
	}
	err = db.Cache.Set(key, b)
	if err != nil {
		logrus.WithError(err).Errorf("Error while saving fake %s", key)
	}
}
