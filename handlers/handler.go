package handlers

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/denouche/go-api-skeleton/middlewares"
	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/dao/fake"
	"github.com/denouche/go-api-skeleton/storage/dao/mock"
	"github.com/denouche/go-api-skeleton/storage/dao/mongodb"
	"github.com/denouche/go-api-skeleton/storage/dao/postgresql"
	"github.com/denouche/go-api-skeleton/storage/validators"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	Mock            bool
	DBInMemory      bool
	DBConnectionURI string
	DBName          string
	Port            int
	LogLevel        string
	LogFormat       string
}

type handlersContext struct {
	db        dao.Database
	validator *validator.Validate
}

func NewRouter(config *Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(gin.Recovery())
	router.Use(middlewares.GetLoggerMiddleware())
	router.Use(middlewares.GetHTTPLoggerMiddleware())

	hc := &handlersContext{}
	if config.Mock {
		hc.db = mock.NewDatabaseMock()
	} else if config.DBInMemory {
		hc.db = fake.NewDatabaseFake()
	} else if strings.HasPrefix(config.DBConnectionURI, "postgresql://") {
		hc.db = postgresql.NewDatabasePostgreSQL(config.DBConnectionURI)
	} else if strings.HasPrefix(config.DBConnectionURI, "mongodb://") {
		hc.db = mongodb.NewDatabaseMongoDB(config.DBConnectionURI, config.DBName)
	} else {
		hc.db = fake.NewDatabaseFake()
	}
	hc.validator = newValidator()

	public := router.Group("/")
	public.Handle(http.MethodGet, "/_health", hc.GetHealth)

	// start: user routes
	public.Handle(http.MethodGet, "/users", hc.GetAllUsers)
	public.Handle(http.MethodPost, "/users", hc.CreateUser)
	public.Handle(http.MethodGet, "/users/:id", hc.GetUser)
	public.Handle(http.MethodPut, "/users/:id", hc.UpdateUser)
	public.Handle(http.MethodDelete, "/users/:id", hc.DeleteUser)
	// end: user routes
	// start: template routes
	public.Handle(http.MethodGet, "/templates", hc.GetAllTemplates)
	public.Handle(http.MethodPost, "/templates", hc.CreateTemplate)
	public.Handle(http.MethodGet, "/templates/:id", hc.GetTemplate)
	public.Handle(http.MethodPut, "/templates/:id", hc.UpdateTemplate)
	public.Handle(http.MethodDelete, "/templates/:id", hc.DeleteTemplate)
	// end: template routes

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

	for k, v := range validators.CustomValidators {
		if v.Validator != nil {
			va.RegisterValidationCtx(k, v.Validator)
		}
	}

	return va
}
