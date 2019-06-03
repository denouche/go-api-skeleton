package dao

import (
	"github.com/denouche/go-api-skeleton/client/model"
)

type Database interface {

	// start: template dao funcs
	GetAllTemplates() ([]*model.Template, error)
	GetTemplateByID(string) (*model.Template, error)
	CreateTemplate(*model.Template) error
	DeleteTemplate(string) error
	UpdateTemplate(template *model.Template) error
	// end: template dao funcs

}
