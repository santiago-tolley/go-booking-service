package clients

import (
	"context"
	"go-booking-service/commons"
	"time"
)

type ClientsService interface {
	Authorize(context.Context, string, string) (string, error)
	Validate(context.Context, string) (string, error)
}

type clientsService struct {
	encoder EncoderDecoder
	users   map[string]string
}

type EncoderDecoder interface {
	Encode(string, string, time.Time) (string, error)
	Decode(string, string) (string, error)
}

func NewClientsServer(e EncoderDecoder, users map[string]string) ClientsService {
	return clientsService{e, users}
}

func (c clientsService) Authorize(ctx context.Context, user, password string) (string, error) {
	if c.users[user] != password {
		return "", ErrInvalidCredentials()
	}

	token, err := c.encoder.Encode(user, commons.JWTSecret, time.Now().Local().Add(commons.JWTExpiration))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c clientsService) Validate(ctx context.Context, token string) (string, error) {
	user, err := c.encoder.Decode(token, commons.JWTSecret)
	if err != nil {
		return "", err
	}
	if _, ok := c.users[user]; !ok {
		return "", ErrUserNotFound()
	}

	return user, err
}
