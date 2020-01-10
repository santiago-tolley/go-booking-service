package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTEncoder struct{}

func (e JWTEncoder) Encode(user, secret string, expiration time.Time) (string, error) {

	if expiration.Before(time.Now()) {
		return "", ErrInvalidDate{}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  expiration.Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err

}

func (e JWTEncoder) Decode(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAlgorithm{token.Header["alg"].(string)}
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user"].(string), nil
	}
	return "", ErrInvalidToken{}
}

type ErrInvalidDate struct{}

func (e ErrInvalidDate) Error() string {
	return "Invalid expiration date"
}

type ErrInvalidAlgorithm struct {
	alg string
}

func (e ErrInvalidAlgorithm) Error() string {
	return fmt.Sprintf("Invalid signing method %v", e.alg)
}

type ErrInvalidToken struct{}

func (e ErrInvalidToken) Error() string {
	return "Invalid JSON web token"
}
