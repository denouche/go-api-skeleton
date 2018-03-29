package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (hc *handlersContext) GetAllUsers(c *gin.Context) {
	users, err := hc.db.GetAllUsers()
	if err != nil {
		utils.JSONErrorWithMessage(c.Writer, model.ErrInternalServer, "Error while getting users")
		return
	}
	utils.JSON(c.Writer, http.StatusOK, users)
}

func (hc *handlersContext) CreateUser(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		logrus.WithError(err).Error("Error while creating user, read data fail")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	userToCreate := model.UserEditable{}
	err = json.Unmarshal(b, &userToCreate)
	if err != nil {
		utils.JSONError(c.Writer, model.ErrBadRequestFormat)
		return
	}

	err = hc.validator.Struct(userToCreate)
	if err != nil {
		utils.JSONError(c.Writer, model.NewDataValidationAPIError(err))
		return
	}

	user := model.User{
		UserEditable: userToCreate,
	}

	err = hc.db.CreateUser(&user)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeDuplicate:
			utils.JSONErrorWithMessage(c.Writer, model.ErrAlreadyExists, "User with the given email already exists")
			return
		default:
			logrus.WithError(err).WithField("type", e.Type).Error("CreateUser: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		logrus.WithError(err).Error("Error while creating user")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusCreated, user)
}
