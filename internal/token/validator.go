package token

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func Validator(tokenString string) (bool, error) {
	signingKey := []byte(os.Getenv("SIGN_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId := claims["sub"]
		if userId == nil {
			return false, fmt.Errorf("user id not found in token")
		}

		if _, err := uuid.Parse(userId.(string)); err != nil {
			return false, fmt.Errorf("invalid user id '%v' in token: %w", userId, err)
		}

		return true, nil
	}

	return false, fmt.Errorf("invalid token")
}
