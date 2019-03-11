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

func (hc *handlersContext) GetAllUsers(c *gin.Context) {
	users, err := hc.db.GetAllUsers()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while getting users")
		utils.JSONErrorWithMessage(c.Writer, model.ErrInternalServer, "Error while getting users")
		return
	}
	utils.JSON(c.Writer, http.StatusOK, users)
}

func (hc *handlersContext) CreateUser(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while creating user, read data fail")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	userToCreate := model.UserEditable{}
	err = json.Unmarshal(b, &userToCreate)
	if err != nil {
		utils.JSONError(c.Writer, model.ErrBadRequestFormat)
		return
	}

	err = hc.validator.StructCtx(c, userToCreate)
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	user := model.User{
		UserEditable: userToCreate,
	}

	err = hc.db.CreateUser(&user)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeDuplicate:
			utils.JSONErrorWithMessage(c.Writer, model.ErrAlreadyExists, "User already exists")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error CreateUser: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while creating user")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusCreated, user)
}

func (hc *handlersContext) GetUser(c *gin.Context) {
	userID := c.Param("id")

	err := hc.validator.VarCtx(c, userID, "required")
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	user, err := hc.db.GetUsersByID(userID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "User not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error GetUser: get user error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get user")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	if user == nil {
		utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "User not found")
		return
	}

	utils.JSON(c.Writer, http.StatusOK, user)
}

func (hc *handlersContext) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	err := hc.validator.VarCtx(c, userID, "required")
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	// check user id given in URL exists
	_, err = hc.db.GetUsersByID(userID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "User to delete not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error DeleteUser: get user error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get user to delete")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	err = hc.db.DeleteUser(userID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "User to delete not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error DeleteUser: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while deleting user")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusNoContent, nil)
}

func (hc *handlersContext) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	err := hc.validator.VarCtx(c, userID, "required")
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	// check user id given in URL exists
	user, err := hc.db.GetUsersByID(userID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "User to update not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("deleteUser: get user error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get user to update")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	// get body and verify data
	b, err := c.GetRawData()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while updating user, read data fail")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	userToUpdate := model.UserEditable{}
	err = json.Unmarshal(b, &userToUpdate)
	if err != nil {
		utils.JSONError(c.Writer, model.ErrBadRequestFormat)
		return
	}

	err = hc.validator.StructCtx(c, userToUpdate)
	if err != nil {
		utils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	user.UserEditable = userToUpdate

	// make the update
	err = hc.db.UpdateUser(user)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			utils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "User to update not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error UpdateUser: Error type not handled")
			utils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while updating user")
		utils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	utils.JSON(c.Writer, http.StatusOK, user)
}
