package fake

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/denouche/go-api-skeleton/utils"

	"github.com/denouche/go-api-skeleton/client/model"
	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/satori/go.uuid"
)

const (
	cacheKeyTemplates = "templates"
)

func (db *DatabaseFake) saveTemplates(templates []*model.Template) {
	data := make([]interface{}, 0)
	for _, v := range templates {
		data = append(data, v)
	}
	db.save(cacheKeyTemplates, data)
}

func (db *DatabaseFake) loadTemplates() []*model.Template {
	templates := make([]*model.Template, 0)
	b, err := db.Cache.Get([]byte(cacheKeyTemplates))
	if err != nil {
		return templates
	}
	err = json.Unmarshal(b, &templates)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Error while unmarshal fake templates")
	}
	return templates
}

func (db *DatabaseFake) GetAllTemplates() ([]*model.Template, error) {
	return db.loadTemplates(), nil
}

func (db *DatabaseFake) GetTemplateByID(templateID string) (*model.Template, error) {
	templates := db.loadTemplates()
	for _, u := range templates {
		if u.ID == templateID {
			return u, nil
		}
	}
	return nil, dao.NewDAOError(dao.ErrTypeNotFound, errors.New("template not found"))
}

func (db *DatabaseFake) CreateTemplate(template *model.Template) error {
	template.ID = uuid.NewV4().String()
	template.CreatedAt = time.Now()

	templates := db.loadTemplates()
	templates = append(templates, template)
	db.saveTemplates(templates)
	return nil
}

func (db *DatabaseFake) DeleteTemplate(templateID string) error {
	templates := db.loadTemplates()
	newTemplates := make([]*model.Template, 0)
	for _, u := range templates {
		if u.ID != templateID {
			newTemplates = append(newTemplates, u)
		}
	}
	db.saveTemplates(newTemplates)
	return nil
}

func (db *DatabaseFake) UpdateTemplate(template *model.Template) error {
	templates := db.loadTemplates()
	var foundTemplate *model.Template
	for _, u := range templates {
		if u.ID == template.ID {
			foundTemplate = u
			break
		}
	}

	if foundTemplate == nil {
		return dao.NewDAOError(dao.ErrTypeNotFound, errors.New("template not found"))
	}

	foundTemplate.TemplateEditable = template.TemplateEditable
	now := time.Now()
	foundTemplate.UpdatedAt = &now
	db.saveTemplates(templates)

	*template = *foundTemplate
	return nil
}
