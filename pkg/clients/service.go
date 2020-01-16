package clients

import (
	"context"
	"go-booking-service/commons"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientsService interface {
	Authorize(context.Context, string, string) (string, error)
	Validate(context.Context, string) (string, error)
	Create(context.Context, string, string) error
}

type clientsService struct {
	encoder EncoderDecoder
	db      *mongo.Database
}

type EncoderDecoder interface {
	Encode(string, string, time.Time) (string, error)
	Decode(string, string) (string, error)
}

func NewClientsServer(e EncoderDecoder, db *mongo.Database) ClientsService {
	return clientsService{e, db}
}

func (c clientsService) Authorize(ctx context.Context, user, password string) (string, error) {
	users := c.db.Collection(commons.MongoClientCollection)
	result := struct {
		User     string
		Password string
	}{}
	filter := bson.D{{"user", user}}
	err := users.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", ErrInvalidCredentials()
	}

	if result.Password != password {
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

	users := c.db.Collection(commons.MongoClientCollection)
	result := struct {
		User     string
		Password string
	}{}
	filter := bson.D{{"user", user}}
	err = users.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", ErrUserNotFound()
	}

	return user, err
}

func (c clientsService) Create(ctx context.Context, user, password string) error {

	users := c.db.Collection(commons.MongoClientCollection)
	result := struct {
		User     string
		Password string
	}{}
	filter := bson.D{{"user", user}}
	err := users.FindOne(context.Background(), filter).Decode(&result)
	if err == nil {
		return ErrUserExists()
	}
	_, err = users.InsertOne(context.Background(), bson.M{user: password})
	if err != nil {
		return err
	}
	// id := res.InsertedID

	return nil
}
