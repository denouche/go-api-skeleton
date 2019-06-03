package fake

import (
	"encoding/json"
	"io/ioutil"

	"github.com/coocood/freecache"
	"github.com/denouche/go-api-skeleton/client/model"
	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/utils"
)

const (
	cacheMaxMemory = 32 * 1024 * 1024 // bytes
)

type DatabaseFake struct {
	Cache *freecache.Cache
}

func NewDatabaseFake(file string) dao.Database {
	result := &DatabaseFake{
		Cache: freecache.NewCache(cacheMaxMemory),
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
	}

	return result
}

func (db *DatabaseFake) save(key string, data []interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		utils.GetLogger().WithError(err).Errorf("Error while marshal fake %s", key)
		_ = db.Cache.Set([]byte(key), []byte("[]"), 0)
		return
	}
	err = db.Cache.Set([]byte(key), b, 0)
	if err != nil {
		utils.GetLogger().WithError(err).Errorf("Error while saving fake %s", key)
	}
}

type Export struct {
	Templates []*model.Template // Template export
}

func (db *DatabaseFake) Export() *Export {
	return &Export{
		Templates: db.loadTemplates(), // Template export
	}
}
