package utilities

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash/fnv"

	"golang.org/x/crypto/bcrypt"
)

// GenerateValueCookie возвращает случайно сгенерированную сроку в base64
func GenerateValueCookie() string  {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	value := base64.URLEncoding.EncodeToString(b)
	return value
}

// HashString хэширует строку с использованием алгоритма FNV-1a
// и возвращает результат в виде шестнадцатеричной строки.
func HashString(s string) string {
	hash := fnv.New64a()
	_, err := hash.Write([]byte(s))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum64())
}

// HashPassword создает хеш пароля.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash сравнивает пароль с его хешем.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}