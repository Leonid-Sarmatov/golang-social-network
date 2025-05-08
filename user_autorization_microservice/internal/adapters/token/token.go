package token

import (
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type tokenJWTAdapter struct {
	secretKey []byte
}

func NewTokenJWTAdapter(secretKey string) *tokenJWTAdapter {
	return &tokenJWTAdapter{
		secretKey: []byte(secretKey),
	}
}

/*
CreateJWTTokenString метод создающий токен для авторизации
*/
func (t *tokenJWTAdapter) CreateToken(data string, minutes int) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": data,
		"nbf":  now.Unix(),
		"exp":  now.Add(time.Duration(minutes) * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/*
ValidateJWTToken метод проверяющий валидность токена
*/
func (t *tokenJWTAdapter) ValidateToken(tokenString string) (string, error) {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := tokenFromString.Claims.(jwt.MapClaims)
	if !ok || !tokenFromString.Valid {
		return "", fmt.Errorf("invalid token")
	} else {
		return claims["name"].(string), nil
	}
}