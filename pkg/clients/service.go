package clients

import (
	"context"
	"go-booking-service/commons"
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errLogger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))

type Service interface {
	Authorize(context.Context, string, string) (string, error)
	Validate(context.Context, string) (string, error)
	Create(context.Context, string, string) error
}

type ClientsService struct {
	encoder EncoderDecoder
	db      *mongo.Database
}

type EncoderDecoder interface {
	Encode(string, string, time.Time) (string, error)
	Decode(string, string) (string, error)
}

type ServiceOption func(*ClientsService)

func WithEncoder(e EncoderDecoder) ServiceOption {
	return func(c *ClientsService) {
		c.encoder = e
	}
}

func WithMongoDB(url, database string) ServiceOption {
	return func(c *ClientsService) {
		db, err := mongo.NewClient(options.Client().ApplyURI(url))
		if err != nil {
			errLogger.Log("message", "could not set up mongo client", "error", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = db.Connect(ctx)
		if err != nil {
			errLogger.Log("message", "could not connect to database", "error", err)
		}
		c.db = db.Database(database)
	}
}

func NewClientsServer(opts ...ServiceOption) *ClientsService {
	c := &ClientsService{}
	for _, options := range opts {
		options(c)
	}
	return c
}

func (c *ClientsService) Authorize(ctx context.Context, user, password string) (string, error) {
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

func (c *ClientsService) Validate(ctx context.Context, token string) (string, error) {
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

func (c *ClientsService) Create(ctx context.Context, user, password string) error {

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
	_, err = users.InsertOne(context.Background(), bson.M{"user": user, "password": password})
	if err != nil {
		return err
	}
	// id := res.InsertedID

	return nil
}
