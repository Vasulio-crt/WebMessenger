package utilities

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateValueCookie() string  {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	value := base64.URLEncoding.EncodeToString(b)
	return value
}