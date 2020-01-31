package server

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AuthorizeEndpoint endpoint.Endpoint
	ValidateEndpoint  endpoint.Endpoint
	BookEndpoint      endpoint.Endpoint
	CheckEndpoint     endpoint.Endpoint
	CreateEndpoint    endpoint.Endpoint
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

func (e Endpoints) Authorize(ctx context.Context, user, password string) (string, error) {
	resp, err := e.AuthorizeEndpoint(ctx, &AuthorizeRequest{User: user, Password: password})
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "authorize failed", "error", err)
		return "", err
	}
	response, ok := resp.(*AuthorizeResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return "", ErrInvalidResponseStructure()
	}
	return response.Token, response.Err
}

func (e Endpoints) Validate(ctx context.Context, token string) (string, error) {
	resp, err := e.ValidateEndpoint(ctx, &ValidateRequest{Token: token})
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "validate failed", "error", err)
		return "", err
	}
	response, ok := resp.(*ValidateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return "", ErrInvalidResponseStructure()
	}
	return response.User, response.Err
}

func (e Endpoints) Create(ctx context.Context, user, password string) error {
	resp, err := e.CreateEndpoint(ctx, &CreateRequest{User: user, Password: password})
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "create failed", "error", err)
		return err
	}
	response, ok := resp.(*CreateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return ErrInvalidResponseStructure()
	}
	return response.Err
}

func MakeEndpoints(p ServerService) Endpoints {
	return Endpoints{
		AuthorizeEndpoint: MakeAuthorizeEndpoint(p),
		ValidateEndpoint:  MakeValidateEndpoint(p),
		BookEndpoint:      MakeBookEndpoint(p),
		CheckEndpoint:     MakeCheckEndpoint(p),
		CreateEndpoint:    MakeCreateEndpoint(p),
	}
}

func MakeAuthorizeEndpoint(p ServerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*AuthorizeRequest)
		if !ok {
			level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
			return &AuthorizeResponse{}, ErrInvalidRequestStructure()
		}
		token, err := p.Authorize(ctx, req.User, req.Password)
		return &AuthorizeResponse{token, err}, nil
	}
}

func MakeValidateEndpoint(p ServerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*ValidateRequest)
		if !ok {
			level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
			return &ValidateResponse{}, ErrInvalidRequestStructure()
		}
		user, err := p.Validate(ctx, req.Token)
		return &ValidateResponse{user, err}, nil
	}
}

func MakeBookEndpoint(p ServerService) endpoint.Endpoint {
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

func MakeCheckEndpoint(p ServerService) endpoint.Endpoint {
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

func MakeCreateEndpoint(p ServerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*CreateRequest)
		if !ok {
			level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
			return &CreateResponse{}, ErrInvalidRequestStructure()
		}
		err := p.Create(ctx, req.User, req.Password)
		return &CreateResponse{err}, nil
	}
}
