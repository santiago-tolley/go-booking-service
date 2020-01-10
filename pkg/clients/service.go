package clients

import (
	"context"
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
		return "", ErrInvalidCredentials{}
	}

	token, err := c.encoder.Encode(user, "very_safe", time.Now().Local().Add(5*time.Minute))

	if err != nil {
		return "", err
	}
	return token, nil
}

func (c clientsService) Validate(ctx context.Context, token string) (string, error) {
	user, err := c.encoder.Decode(token, "very_safe")

	if _, ok := c.users[user]; !ok {
		return "", ErrUserNotFound{user}
	}

	return user, err
}
