package clients

import (
	"context"
	"os"
	"time"

	"go-booking-service/commons"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))

func init() {
	logger = level.NewFilter(logger, commons.LoggingLevel)
	logger = kitlog.With(logger, "origin", "Clients", "caller", kitlog.DefaultCaller)
}

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
	level.Info(logger).Log("message", "Setting up encoder")
	return func(c *ClientsService) {
		c.encoder = e
	}
}

func WithMongoDB(url, database string) ServiceOption {
	return func(c *ClientsService) {
		level.Info(logger).Log("message", "Setting up mongodb database")
		db, err := mongo.NewClient(options.Client().ApplyURI(url))
		if err != nil {
			level.Error(logger).Log("message", "could not set up mongo client", "error", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = db.Connect(ctx)
		if err != nil {
			level.Error(logger).Log("message", "could not connect to database", "error", err)
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
	level.Info(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Authorize attempt", "user", user, "password", password)
	users := c.db.Collection(commons.MongoClientCollection)
	result := struct {
		User     string
		Password string
	}{}
	filter := bson.D{{"user", user}}
	err := users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "error retrieving user", "error", err)
		return "", ErrInvalidCredentials()
	}

	if result.Password != password {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "error password does not match")
		return "", ErrInvalidCredentials()
	}

	token, err := c.encoder.Encode(user, commons.JWTSecret, time.Now().Local().Add(commons.JWTExpiration))
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "could not generate token", "error", err)
		return "", err
	}
	return token, nil
}

func (c *ClientsService) Validate(ctx context.Context, token string) (string, error) {
	level.Info(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Validate attempt", "token", token)
	user, err := c.encoder.Decode(token, commons.JWTSecret)
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid token", "error", err)
		return "", err
	}

	users := c.db.Collection(commons.MongoClientCollection)
	result := struct {
		User     string
		Password string
	}{}
	filter := bson.D{{"user", user}}
	err = users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "error retrieving user", "error", err)
		return "", ErrUserNotFound()
	}

	return user, err
}

func (c *ClientsService) Create(ctx context.Context, user, password string) error {
	level.Info(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Create attempt", "user", user, "password", password)
	users := c.db.Collection(commons.MongoClientCollection)

	_, err := users.InsertOne(ctx, bson.M{"user": user, "password": password})
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "error inserting in database", "error", err)
		return ErrUserExists()
	}

	return nil
}
