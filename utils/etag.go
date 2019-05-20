package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	HeaderNameIfMatch = "If-Match"
)

func getHash(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

func GenerateEtag(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	str := string(b)
	tag := fmt.Sprintf("\"%d-%s\"", len(str), getHash(str))
	return tag, nil
}

func IsSameVersion(c *gin.Context, resource interface{}) bool {
	etag, err := GenerateEtag(resource)
	if err != nil {
		return false
	}

	ifMatch := c.GetHeader(HeaderNameIfMatch)
	if ifMatch == "" {
		return false
	}

	if ifMatch != etag {
		return false
	}

	return true
}
