package rooms

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	BookEndpoint  endpoint.Endpoint
	CheckEndpoint endpoint.Endpoint
}

func (e Endpoints) Book(ctx context.Context, token string, date time.Time) (int, error) {
	resp, err := e.BookEndpoint(ctx, BookRequest{Token: token, Date: date})
	if err != nil {
		return 0, err
	}
	response, ok := resp.(BookResponse)
	if !ok {
		return 0, ErrInvalidResponseStructure{}
	}

	return response.Id, response.Err
}

func (e Endpoints) Check(ctx context.Context, date time.Time) (int, error) {
	resp, err := e.CheckEndpoint(ctx, CheckRequest{Date: date})
	if err != nil {
		return 0, err
	}
	response, ok := resp.(CheckResponse)
	if !ok {
		return 0, ErrInvalidResponseStructure{}
	}

	return response.Available, response.Err
}

func MakeEndpoints(p RoomsService) Endpoints {
	return Endpoints{
		BookEndpoint:  MakeBookEndpoint(p),
		CheckEndpoint: MakeCheckEndpoint(p),
	}
}

func MakeBookEndpoint(p RoomsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(BookRequest)
		if !ok {
			return BookResponse{}, ErrInvalidRequestStructure{}
		}
		id, err := p.Book(ctx, req.Token, req.Date)

		return BookResponse{id, err}, nil
	}
}

func MakeCheckEndpoint(p RoomsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CheckRequest)
		if !ok {
			return CheckResponse{}, ErrInvalidRequestStructure{}
		}
		available, err := p.Check(ctx, req.Date)
		if err != nil {
			return CheckResponse{0, err}, nil
		}
		return CheckResponse{available, nil}, nil
	}
}
