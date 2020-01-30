package rooms

import (
	"context"
	"go-booking-service/commons"
	"go-booking-service/pb"

	"github.com/go-kit/kit/log/level"
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

func encodeGRPCBookRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*BookRequest)
	level.Info(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "Book attempt")
	if !ok {
		level.Error(logger).Log("message", "encoding: invalid book request structure")
		return &pb.BookRequest{}, ErrInvalidRequestStructure()
	}

	var correlationId string
	correlationId, ok = ctx.Value(commons.ContextKeyCorrelationID).(string)
	if !ok {
		level.Error(logger).Log("message", "encoding: no correlation id in context")
		return &pb.BookRequest{}, ErrNoCorrelationId()
	}

	return &pb.BookRequest{
		Token:         req.Token,
		Date:          req.Date.Unix(),
		CorrelationId: correlationId,
	}, nil
}

func decodeGRPCBookResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.BookResponse)
	if !ok {
		level.Error(logger).Log("message", "encoding: invalid book response structure")
		return &BookResponse{}, ErrInvalidResponseStructure()
	}

	// ctx = context.WithValue(ctx, commons.ContextKeyCorrelationID, reply.CorrelationId)

	return &BookResponse{
		Id:  int(reply.Id),
		Err: str2err(reply.Error),
	}, nil
}

func encodeGRPCCheckRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*CheckRequest)
	if !ok {
		return &pb.CheckRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.CheckRequest{
		Date: req.Date.Unix(),
	}, nil
}

func decodeGRPCCheckResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.CheckResponse)
	if !ok {
		return &CheckResponse{}, ErrInvalidResponseStructure()
	}
	return &CheckResponse{
		Available: int(reply.Available),
		Err:       str2err(reply.Error),
	}, nil
}

func str2err(s string) error {
	switch s {
	case "":
		return nil
	case InvalidRequestStructure:
		return ErrInvalidRequestStructure()
	case InvalidResponseStructure:
		return ErrInvalidResponseStructure()
	case NoRoomAvailable:
		return ErrNoRoomAvailable()
	default:
		return ErrorWithMsg{s}
	}
}
