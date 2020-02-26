package server

import (
	"context"
	"os"
	"time"

	"go-booking-service/commons"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))

func init() {
	logger = level.NewFilter(logger, commons.LoggingLevel)
	logger = kitlog.With(logger, "origin", "Server", "caller", kitlog.DefaultCaller)
}

type ClientsService interface {
	Authorize(context.Context, string, string) (string, error)
	Validate(context.Context, string) (string, error)
	Create(context.Context, string, string) error
}

type RoomService interface {
	Book(context.Context, string, time.Time) (int, error)
	Check(context.Context, time.Time) (int, error)
}

func NewServer(clientsClient ClientsService, roomsClient RoomService) ServerService {
	return ServerService{ClientsClient: clientsClient, RoomClient: roomsClient}
}

type ServerService struct {
	ClientsClient ClientsService
	RoomClient    RoomService
}

func (p ServerService) Authorize(ctx context.Context, user, password string) (string, error) {
	level.Info(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Authorize attempt", "user", user, "password", password)
	token, err := p.ClientsClient.Authorize(ctx, user, password)
	return token, err
}

func (p ServerService) Validate(ctx context.Context, token string) (string, error) {
	level.Info(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Validate attempt", "token", token)
	user, err := p.ClientsClient.Validate(ctx, token)
	return user, err
}

func (p ServerService) Create(ctx context.Context, user, password string) error {
	level.Info(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Create attempt", "user", user, "password", password)
	err := p.ClientsClient.Create(ctx, user, password)
	return err
}

func (p ServerService) Book(ctx context.Context, token string, date time.Time) (int, error) {
	level.Info(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Book attempt", "token", token, "date", date)
	book_id, err := p.RoomClient.Book(ctx, token, date)
	return book_id, err
}

func (p ServerService) Check(ctx context.Context, date time.Time) (int, error) {
	level.Info(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Check attempt", "date", date)
	available, err := p.RoomClient.Check(ctx, date)
	return available, err
}
