package fake

import (
	"encoding/json"
	"time"

	"github.com/denouche/go-api-skeleton/utils"

	"github.com/allegro/bigcache"
	"github.com/denouche/go-api-skeleton/storage/dao"
)

type DatabaseFake struct {
	Cache *bigcache.BigCache
}

func NewDatabaseFake() dao.Database {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	if err != nil {
		utils.GetLogger().WithError(err).Fatal("Error while instantiate cache")
	}
	return &DatabaseFake{
		Cache: cache,
	}
}

func (db *DatabaseFake) save(key string, data []interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		utils.GetLogger().WithError(err).Errorf("Error while marshal fake %s", key)
		db.Cache.Set(key, []byte("[]"))
		return
	}
	err = db.Cache.Set(key, b)
	if err != nil {
		utils.GetLogger().WithError(err).Errorf("Error while saving fake %s", key)
	}
}
