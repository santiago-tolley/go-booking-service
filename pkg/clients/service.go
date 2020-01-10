package clients

import (
	"context"
	"time"
)

type ClientsService struct {
	encoder EncoderDecoder
	users   map[string]string
}

type EncoderDecoder interface {
	Encode(string, string, time.Time) (string, error)
	Decode(string, string) (string, error)
}

func NewClientsServer(e EncoderDecoder, users map[string]string) ClientsService {
	return ClientsService{e, users}
}

func (c ClientsService) Authorize(ctx context.Context, user, password string) (string, error) {
	if c.users[user] != password {
		return "", ErrInvalidCredentials{}
	}

	token, err := c.encoder.Encode(user, "very_safe", time.Now().Local().Add(5*time.Minute))

	if err != nil {
		return "", err
	}
	return token, nil
}

func (c ClientsService) Validate(ctx context.Context, token string) (string, error) {
	user, err := c.encoder.Decode(token, "very_safe")

	if _, ok := c.users[user]; !ok {
		return "", ErrUserNotFound{user}
	}

	return user, err
}
