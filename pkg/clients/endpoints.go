package clients

import (
	"context"
	"go-booking-service/commons"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
)

type Endpoints struct {
	AuthorizeEndpoint endpoint.Endpoint
	ValidateEndpoint  endpoint.Endpoint
	CreateEndpoint    endpoint.Endpoint
}

func (e Endpoints) Authorize(ctx context.Context, user, password string) (string, error) {
	level.Info(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "in clients > endpoints > authorize")
	resp, err := e.AuthorizeEndpoint(ctx, &AuthorizeRequest{User: user, Password: password})
	if err != nil {
		return "", err
	}
	response, ok := resp.(*AuthorizeResponse)
	if !ok {
		return "", ErrInvalidResponseStructure()
	}

	return response.Token, response.Err
}

func (e Endpoints) Validate(ctx context.Context, token string) (string, error) {
	resp, err := e.ValidateEndpoint(ctx, &ValidateRequest{Token: token})
	if err != nil {
		return "", err
	}
	response, ok := resp.(*ValidateResponse)
	if !ok {
		return "", ErrInvalidResponseStructure()
	}

	return response.User, response.Err
}

func (e Endpoints) Create(ctx context.Context, user, password string) error {
	resp, err := e.CreateEndpoint(ctx, &CreateRequest{User: user, Password: password})
	if err != nil {
		return err
	}
	response, ok := resp.(*CreateResponse)
	if !ok {
		return ErrInvalidResponseStructure()
	}

	return response.Err
}

func MakeEndpoints(c Service) Endpoints {
	return Endpoints{
		AuthorizeEndpoint: MakeAuthorizeEndpoint(c),
		ValidateEndpoint:  MakeValidateEndpoint(c),
		CreateEndpoint:    MakeCreateEndpoint(c),
	}
}

func MakeAuthorizeEndpoint(c Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*AuthorizeRequest)
		if !ok {
			return &AuthorizeResponse{}, ErrInvalidRequestStructure()
		}

		level.Debug(logger).Log("correlation ID", req.CorrelationID, "message", "in make endpoint result authorize")
		ctx = context.WithValue(ctx, commons.ContextKeyCorrelationID, req.CorrelationID)
		token, err := c.Authorize(ctx, req.User, req.Password)
		return &AuthorizeResponse{token, err}, nil
	}
}

func MakeValidateEndpoint(c Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*ValidateRequest)
		if !ok {
			return &ValidateResponse{}, ErrInvalidRequestStructure()
		}

		user, err := c.Validate(ctx, req.Token)
		return &ValidateResponse{user, err}, nil
	}
}

func MakeCreateEndpoint(c Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*CreateRequest)
		if !ok {
			return &CreateResponse{}, ErrInvalidRequestStructure()
		}

		err := c.Create(ctx, req.User, req.Password)
		return &CreateResponse{err}, nil
	}
}
