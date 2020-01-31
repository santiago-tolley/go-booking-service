package rooms

import (
	"context"
	"go-booking-service/commons"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
)

type Endpoints struct {
	BookEndpoint  endpoint.Endpoint
	CheckEndpoint endpoint.Endpoint
}

func (e Endpoints) Book(ctx context.Context, token string, date time.Time) (int, error) {
	resp, err := e.BookEndpoint(ctx, &BookRequest{Token: token, Date: date})
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "book failed", "error", err)
		return 0, err
	}
	response, ok := resp.(*BookResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return 0, ErrInvalidResponseStructure()
	}

	return response.Id, response.Err
}

func (e Endpoints) Check(ctx context.Context, date time.Time) (int, error) {
	resp, err := e.CheckEndpoint(ctx, &CheckRequest{Date: date})
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "check failed", "error", err)
		return 0, err
	}
	response, ok := resp.(*CheckResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return 0, ErrInvalidResponseStructure()
	}

	return response.Available, response.Err
}

func MakeEndpoints(p Service) Endpoints {
	return Endpoints{
		BookEndpoint:  MakeBookEndpoint(p),
		CheckEndpoint: MakeCheckEndpoint(p),
	}
}

func MakeBookEndpoint(p Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*BookRequest)
		if !ok {
			level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
			return &BookResponse{}, ErrInvalidRequestStructure()
		}
		id, err := p.Book(ctx, req.Token, req.Date)

		return &BookResponse{id, err}, nil
	}
}

func MakeCheckEndpoint(p Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*CheckRequest)
		if !ok {
			level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
			return &CheckResponse{}, ErrInvalidRequestStructure()
		}
		available, err := p.Check(ctx, req.Date)

		return &CheckResponse{available, err}, nil
	}
}
