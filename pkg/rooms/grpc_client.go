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
	correlationId, ok := ctx.Value(commons.ContextKeyCorrelationID).(uuid.UUID)
	if !ok {
		level.Error(logger).Log("message", "setting metadata", "error", "no correlation id in context")
		return ctx
	}

	(*md)[commons.MetaDataKeyCorrelationID] = append((*md)[commons.MetaDataKeyCorrelationID], correlationId.String())
	return ctx
}

func encodeGRPCBookRequest(ctx context.Context, request interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding book request")
	req, ok := request.(*BookRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding: invalid book request structure")
		return &pb.BookRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.BookRequest{
		Token: req.Token,
		Date:  req.Date.Unix(),
	}, nil
}

func decodeGRPCBookResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding book response")
	reply, ok := grpcReply.(*pb.BookResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding: invalid book response structure")
		return &BookResponse{}, ErrInvalidResponseStructure()
	}
	return &BookResponse{
		Id:  int(reply.Id),
		Err: str2err(reply.Error),
	}, nil
}

func encodeGRPCCheckRequest(ctx context.Context, request interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding check request")
	req, ok := request.(*CheckRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding: invalid check request structure")
		return &pb.CheckRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.CheckRequest{
		Date: req.Date.Unix(),
	}, nil
}

func decodeGRPCCheckResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding check response")
	reply, ok := grpcReply.(*pb.CheckResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding: invalid check response structure")
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
