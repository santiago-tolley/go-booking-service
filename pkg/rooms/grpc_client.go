package rooms

import (
	"context"

	"go-booking-service/commons"
	"go-booking-service/pb"

	"github.com/go-kit/kit/log/level"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGRPCClient(conn *grpc.ClientConn) Endpoints {
	bookEndpoint := grpctransport.NewClient(
		conn,
		"pb.Rooms",
		"Book",
		encodeGRPCBookRequest,
		decodeGRPCBookResponse,
		pb.BookResponse{},
		grpctransport.ClientBefore(setMDCorrelationID),
	).Endpoint()

	checkEndpoint := grpctransport.NewClient(
		conn,
		"pb.Rooms",
		"Check",
		encodeGRPCCheckRequest,
		decodeGRPCCheckResponse,
		pb.CheckResponse{},
		grpctransport.ClientBefore(setMDCorrelationID),
	).Endpoint()

	return Endpoints{
		BookEndpoint:  bookEndpoint,
		CheckEndpoint: checkEndpoint,
	}
}

func setMDCorrelationID(ctx context.Context, md *metadata.MD) context.Context {
	level.Debug(logger).Log("message", "setting metadata", "context", ctx.Value(commons.ContextKeyCorrelationID))
	correlationId, ok := ctx.Value(commons.ContextKeyCorrelationID).(uuid.UUID)
	if !ok {
		level.Error(logger).Log("message", "setting metadata", "error", "no correlation id in context")
		return ctx
	}

	level.Debug(logger).Log("message", "setting metadata", "context", correlationId.String())
	(*md)[commons.MetaDataKeyCorrelationID] = append((*md)[commons.MetaDataKeyCorrelationID], correlationId.String())
	return ctx
}

func encodeGRPCBookRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*BookRequest)
	if !ok {
		level.Error(logger).Log("message", "encoding: invalid book request structure")
		return &pb.BookRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.BookRequest{
		Token: req.Token,
		Date:  req.Date.Unix(),
	}, nil
}

func decodeGRPCBookResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.BookResponse)
	if !ok {
		level.Error(logger).Log("message", "encoding: invalid book response structure")
		return &BookResponse{}, ErrInvalidResponseStructure()
	}
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
