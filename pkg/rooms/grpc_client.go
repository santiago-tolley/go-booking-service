package rooms

import (
	"context"
	"errors"
	"go-booking-service/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn) Endpoints {
	bookEndpoint := grpctransport.NewClient(
		conn,
		"pb.Rooms",
		"Book",
		encodeGRPCBookRequest,
		decodeGRPCBookResponse,
		pb.BookResponse{},
	).Endpoint()

	checkEndpoint := grpctransport.NewClient(
		conn,
		"pb.Rooms",
		"Check",
		encodeGRPCCheckRequest,
		decodeGRPCCheckResponse,
		pb.CheckResponse{},
	).Endpoint()

	return Endpoints{
		BookEndpoint:  bookEndpoint,
		CheckEndpoint: checkEndpoint,
	}
}

func encodeGRPCBookRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(BookRequest)
	if !ok {
		return &pb.BookRequest{}, ErrInvalidRequestStructure{}
	}
	return &pb.BookRequest{
		Token: req.Token,
		Date:  req.Date.Unix(),
	}, nil
}

func decodeGRPCBookResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.BookResponse)
	if !ok {
		return BookResponse{}, ErrInvalidResponseStructure{}
	}
	return BookResponse{
		Id:  int(reply.Id),
		Err: str2err(reply.Error),
	}, nil
}

func encodeGRPCCheckRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(CheckRequest)
	if !ok {
		return &pb.CheckRequest{}, ErrInvalidRequestStructure{}
	}
	return &pb.CheckRequest{
		Date: req.Date.Unix(),
	}, nil
}

func decodeGRPCCheckResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.CheckResponse)
	if !ok {
		return CheckResponse{}, ErrInvalidResponseStructure{}
	}
	return CheckResponse{
		Available: int(reply.Available),
		Err:       str2err(reply.Error),
	}, nil
}

func str2err(s string) error { // todo
	if s == "" {
		return nil
	}
	return errors.New(s)
}
