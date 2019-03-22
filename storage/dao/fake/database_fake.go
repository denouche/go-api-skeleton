package fake

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/allegro/bigcache"
	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/denouche/go-api-skeleton/utils"
)

type DatabaseFake struct {
	Cache *bigcache.BigCache
}

func NewDatabaseFake(file string) dao.Database {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	if err != nil {
		utils.GetLogger().WithError(err).Fatal("Error while instantiate cache")
	}

	result := &DatabaseFake{
		Cache: cache,
	}

	if file != "" {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			utils.GetLogger().WithError(err).Error("error while loading data for in memory database")
		}

		var export Export
		err = json.Unmarshal(data, &export)
		if err != nil {
			utils.GetLogger().WithError(err).Error("error while reading data from file for in memory database")
		}

		result.saveTemplates(export.Templates) // Template export
		result.saveUsers(export.Users)
	}

	return result
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

type Export struct {
	Templates []*model.Template // Template export
	Users     []*model.User
}

func (db *DatabaseFake) Export() *Export {
	return &Export{
		Templates: db.loadTemplates(), // Template export
		Users:     db.loadUsers(),
	}
}
