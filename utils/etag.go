package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
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

func IsSameVersion(expectedEtag string, resource interface{}) bool {
	etag, err := GenerateEtag(resource)
	if err != nil {
		return false
	}

	if expectedEtag == "" {
		return false
	}

	if expectedEtag != etag {
		return false
	}

	return true
}
