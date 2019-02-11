package fake

import (
	"encoding/json"
	"time"

	"github.com/allegro/bigcache"
	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/utils"
)

type DatabaseFake struct {
	Cache *bigcache.BigCache
}

func NewDatabaseFake() dao.Database {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	if err != nil {
		utils.GetLogger(nil).Errorw("Error while instantiate cache",
			"error", err)
	}
	return &DatabaseFake{
		Cache: cache,
	}
}

func (db *DatabaseFake) save(key string, data []interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		utils.GetLogger(nil).Errorw("Error while marshal fake",
			"key", key,
			"error", err)
		db.Cache.Set(key, []byte("[]"))
		return
	}
	err = db.Cache.Set(key, b)
	if err != nil {
		utils.GetLogger(nil).Errorw("Error while saving fake",
			"key", key,
			"error", err)
	}
}
