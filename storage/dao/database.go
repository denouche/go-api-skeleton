package dao

import (
	"github.com/denouche/go-api-skeleton/storage/model"
)

type Database interface {

	// start: user dao funcs
	GetAllUsers() ([]*model.User, error)
	GetUserByID(string) (*model.User, error)
	CreateUser(*model.User) error
	DeleteUser(string) error
	UpdateUser(*model.User) error
	// end: user dao funcs
	// start: template dao funcs
	GetAllTemplates() ([]*model.Template, error)
	GetTemplateByID(string) (*model.Template, error)
	CreateTemplate(*model.Template) error
	DeleteTemplate(string) error
	UpdateTemplate(template *model.Template) error
	// end: template dao funcs

}
