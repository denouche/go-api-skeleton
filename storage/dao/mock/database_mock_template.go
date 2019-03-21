package mock

import (
	"github.com/denouche/go-api-skeleton/storage/model"
)

func (db *DatabaseMock) GetAllTemplates() ([]*model.Template, error) {
	args := db.Called()
	return args.Get(0).([]*model.Template), args.Error(1)
}

func (db *DatabaseMock) GetTemplateByID(id string) (*model.Template, error) {
	args := db.Called()
	return args.Get(0).(*model.Template), args.Error(1)
}

func (db *DatabaseMock) CreateTemplate(template *model.Template) error {
	args := db.Called()
	return args.Error(0)
}

func (db *DatabaseMock) DeleteTemplate(id string) error {
	args := db.Called()
	return args.Error(0)
}

func (db *DatabaseMock) UpdateTemplate(template *model.Template) error {
	args := db.Called()
	return args.Error(0)
}
