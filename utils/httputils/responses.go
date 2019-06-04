package httputils

import (
	"encoding/json"
	"net/http"

	"github.com/denouche/go-api-skeleton/pkg/client/model"
	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

func JSONOK(c *gin.Context, data interface{}) {
	if utils.IsSameVersion(c.GetHeader(HeaderNameIfNoneMatch), data) {
		JSON(c.Writer, http.StatusNotModified, data)
		return
	}
	JSON(c.Writer, http.StatusOK, data)
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set(HeaderNameContentType, HeaderValueApplicationJSONUTF8)
	w.WriteHeader(status)
	if data != nil {
		etag, err := utils.GenerateEtag(data)
		if err == nil {
			w.Header().Set(HeaderNameAccessControlExposeHeaders, HeaderNameETag)
			w.Header().Set(HeaderNameETag, etag)
		}
		json.NewEncoder(w).Encode(data)
	}
}

func JSONError(w http.ResponseWriter, e model.APIError) {
	if e.Headers != nil {
		for k, headers := range e.Headers {
			for _, headerValue := range headers {
				w.Header().Add(k, headerValue)
			}
		}
	}
	JSON(w, e.HTTPCode, e)
}

func JSONErrorWithMessage(w http.ResponseWriter, e model.APIError, message string) {
	e.Description = message
	JSONError(w, e)
}

func YAML(w http.ResponseWriter, status int, data string) {
	w.Header().Set(HeaderNameContentType, HeaderValueApplicationYAML)
	w.WriteHeader(status)
	w.Write([]byte(data))
}
