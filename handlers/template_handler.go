package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/denouche/go-api-skeleton/storage/validators"
	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

// @openapi:path
// /templates:
//	get:
//		description: "Get all the templates"
//		responses:
//			200:
//				description: "The array containing the templates"
//				content:
//					application/json:
//						schema:
//							type: "array"
//							items:
//								$ref: "#/components/schemas/Template"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) GetAllTemplates(c *gin.Context) {
	templates, err := hc.db.GetAllTemplates()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while getting templates")
		utils.JSONErrorWithMessage(c.Writer, model.ErrInternalServer, "Error while getting templates")
		return
	}
	utils.JSON(c.Writer, http.StatusOK, templates)
}

// @openapi:path
// /templates:
//	post:
//		description: "Create a new template"
//		requestBody:
//			description: The template data.
//			required: true
//			content:
//				application/json:
//					schema:
//						$ref: "#/components/schemas/TemplateEditable"
//		responses:
//			201:
//				description: "The created template"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/Template"
//			400:
//				description: "This error occurs when the request is not correct (bad body format, validation error)"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			409:
//				description: "This error occurs when the new entity is in conflict with exiting one (duplicated)"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) CreateTemplate(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while creating template, read data fail")
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
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error CreateTemplate: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while creating template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusCreated, template)
}

// @openapi:path
// /templates/{templateID}:
//	get:
//		description: "Get a template"
//		parameters:
//		- in: path
//		  name: templateID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The template id to get"
//		responses:
//			200:
//				description: "The templates with id `templateID`"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/Template"
//			404:
//				description: "Template not found"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "required")
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
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error GetTemplate: get template error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	if template == nil {
		utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template not found")
		return
	}

	utils.JSON(c.Writer, http.StatusOK, template)
}

// @openapi:path
// /templates/{templateID}:
//	delete:
//		description: "Delete a template"
//		parameters:
//		- in: path
//		  name: templateID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The template id to delete"
//		responses:
//			204:
//				description: "Templates with id `templateID` deleted"
//			404:
//				description: "Template not found"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "required")
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
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error DeleteTemplate: get template error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get template to delete")
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
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error DeleteTemplate: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while deleting template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusNoContent, nil)
}

// @openapi:path
// /templates/{templateID}:
//	put:
//		description: "Update a template"
//		parameters:
//		- in: path
//		  name: templateID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The template id to update"
//		requestBody:
//			description: The template data.
//			required: true
//			content:
//				application/json:
//					schema:
//						$ref: "#/components/schemas/TemplateEditable"
//		responses:
//			201:
//				description: "The updated template"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/Template"
//			400:
//				description: "This error occurs when the request is not correct (bad body format, validation error)"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			404:
//				description: "Template not found"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "required")
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
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("UpdateTemplate: get template error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get template to update")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	// get body and verify data
	b, err := c.GetRawData()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while updating template, read data fail")
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
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error UpdateTemplate: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while updating template")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusOK, template)
}
