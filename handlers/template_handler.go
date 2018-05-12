package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/denouche/go-api-skeleton/storage/validators"
	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (hc *handlersContext) GetAllTemplates(c *gin.Context) {
	templates, err := hc.db.GetAllTemplates()
	if err != nil {
		logrus.WithError(err).Error("error while getting templates")
		utils.JSONErrorWithMessage(c.Writer, model.ErrInternalServer, "Error while getting templates")
		return
	}
	utils.JSON(c.Writer, http.StatusOK, templates)
}

func (hc *handlersContext) CreateTemplate(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		logrus.WithError(err).Error("error while creating template, read data fail")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	templateToCreate := model.TemplateEditable{}
	err = json.Unmarshal(b, &templateToCreate)
	if err != nil {
		utils.JSONError(c.Writer, model.ErrBadRequestFormat)
		return
	}

	err = hc.validator.StructCtx(c, templateToCreate)
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	template := model.Template{
		TemplateEditable: templateToCreate,
	}

	err = hc.db.CreateTemplate(&template)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeDuplicate:
			utils.JSONErrorWithMessage(c.Writer, model.ErrAlreadyExists, "Template already exists")
			return
		default:
			logrus.WithError(err).WithField("type", e.Type).Error("error CreateTemplate: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("error while creating template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusCreated, template)
}

func (hc *handlersContext) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "uuid4")
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	template, err := hc.db.GetTemplatesByID(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template not found")
			return
		default:
			logrus.WithError(err).WithField("type", e.Type).Error("error GetTemplate: get template error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("error while get template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	if template == nil {
		utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template not found")
		return
	}

	utils.JSON(c.Writer, http.StatusOK, template)
}

func (hc *handlersContext) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "uuid4")
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	// check template id given in URL exists
	_, err = hc.db.GetTemplatesByID(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to delete not found")
			return
		default:
			logrus.WithError(err).WithField("type", e.Type).Error("error DeleteTemplate: get template error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("error while get template to delete")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	err = hc.db.DeleteTemplate(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to delete not found")
			return
		default:
			logrus.WithError(err).WithField("type", e.Type).Error("error DeleteTemplate: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("error while deleting template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusNoContent, nil)
}

func (hc *handlersContext) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "uuid4")
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	// check template id given in URL exists
	template, err := hc.db.GetTemplatesByID(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to update not found")
			return
		default:
			logrus.WithError(err).WithField("type", e.Type).Error("deleteTemplate: get template error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("error while get template to update")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	// get body and verify data
	b, err := c.GetRawData()
	if err != nil {
		logrus.WithError(err).Error("error while updating template, read data fail")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	templateToUpdate := model.TemplateEditable{}
	err = json.Unmarshal(b, &templateToUpdate)
	if err != nil {
		utils.JSONError(c.Writer, model.ErrBadRequestFormat)
		return
	}

	err = hc.validator.StructCtx(c, templateToUpdate)
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	template.TemplateEditable = templateToUpdate

	// make the update
	err = hc.db.UpdateTemplate(template)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to update not found")
			return
		default:
			logrus.WithError(err).WithField("type", e.Type).Error("error UpdateTemplate: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("error while deleting template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusOK, template)
}
