package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTEncoder struct{}

func (e JWTEncoder) Encode(user, secret string, expiration time.Time) (string, error) {

	if expiration.Before(time.Now()) {
		return "", ErrInvalidDate()
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
			return nil, ErrInvalidAlgorithm()
		}
		return []byte(secret), nil
	})

	if err != nil {
		if err.Error() == "Token is expired" {
			return "", ErrExpiredToken()
		}
		return "", ErrInvalidToken()
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user"].(string), nil
	}

	return "", ErrInvalidToken()
}
