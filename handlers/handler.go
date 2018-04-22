package handlers

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	Mock            bool
	DBConnectionURI string
	Port            int
	LogLevel        string
	LogFormat       string
}

type handlersContext struct {
	db        dao.Database
	validator *validator.Validate
}

func NewRouter(config *Config) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.HandleMethodNotAllowed = true

	hc := &handlersContext{}
	if config.Mock {
		hc.db = dao.NewDatabaseFake()
	} else {
		hc.db = dao.NewDatabasePostgreSQL(config.DBConnectionURI)
	}
	hc.validator = newValidator()

	public := router.Group("/")
	public.Handle(http.MethodGet, "/_health", hc.GetHealth)

	public.Handle(http.MethodGet, "/users", hc.GetAllUsers)
	public.Handle(http.MethodPost, "/users", hc.CreateUser)
	public.Handle(http.MethodGet, "/users/:id", hc.GetUser)
	public.Handle(http.MethodPut, "/users/:id", hc.UpdateUser)
	public.Handle(http.MethodDelete, "/users/:id", hc.DeleteUser)

	return router
}

func newValidator() *validator.Validate {
	va := validator.New()

	va.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)
		if len(name) < 1 {
			return ""
		}
		return name[0]
	})

	for k, v := range model.CustomValidators {
		if v.Validator != nil {
			va.RegisterValidation(k, v.Validator)
		}
	}

	return va
}
