package util

import (
	"aherman/src/enums"

	"github.com/golang-jwt/jwt"
)

// JWTToString todo
func JWTToString(t *jwt.Token) (string, error) {
	tokenString, err := t.SignedString([]byte(enums.AppSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
