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
// /users:
//	get:
//		description: "Get all the users"
//		responses:
//			200:
//				description: "The array containing the users"
//				content:
//					application/json:
//						schema:
//							type: "array"
//							items:
//								$ref: "#/components/schemas/User"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *handlersContext) GetAllUsers(c *gin.Context) {
	users, err := hc.db.GetAllUsers()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while getting users")
		utils.JSONErrorWithMessage(c.Writer, model.ErrInternalServer, "Error while getting users")
		return
	}
	utils.JSON(c.Writer, http.StatusOK, users)
}

// @openapi:path
// /users:
//	post:
//		description: "Create a new user"
//		requestBody:
//			description: The user data.
//			required: true
//			content:
//				application/vnd.api+json:
//					schema:
//						$ref: "#/components/schemas/UserEditable"
//		responses:
//			201:
//				description: "The array containing the users"
//				content:
//					application/json:
//						schema:
//							type: "array"
//							items:
//								$ref: "#/components/schemas/User"
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

// @openapi:path
// /users/{userID}:
//	get:
//		description: "Get a user"
//		parameters:
//		- in: path
//		  name: userID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The user id to get"
//		responses:
//			200:
//				description: "The users with id `userID`"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/User"
//			404:
//				description: "User not found"
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

// @openapi:path
// /users/{userID}:
//	delete:
//		description: "Delete a user"
//		parameters:
//		- in: path
//		  name: userID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The user id to delete"
//		responses:
//			204:
//				description: "Users with id `userID` deleted"
//			404:
//				description: "User not found"
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

// @openapi:path
// /users/{userID}:
//	put:
//		description: "Update a user"
//		parameters:
//		- in: path
//		  name: userID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The user id to update"
//		requestBody:
//			description: The user data.
//			required: true
//			content:
//				application/vnd.api+json:
//					schema:
//						$ref: "#/components/schemas/UserEditable"
//		responses:
//			201:
//				description: "The array containing the users"
//				content:
//					application/json:
//						schema:
//							type: "array"
//							items:
//								$ref: "#/components/schemas/User"
//			400:
//				description: "This error occurs when the request is not correct (bad body format, validation error)"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			404:
//				description: "User not found"
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
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("UpdateUser: get user error type not handled")
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
