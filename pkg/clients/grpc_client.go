package clients

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log/level"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go-booking-service/commons"
	"go-booking-service/pb"
)

func NewGRPCClient(conn *grpc.ClientConn) Endpoints {
	authorizeEndpoint := grpctransport.NewClient(
		conn,
		"pb.Clients",
		"Authorize",
		encodeGRPCAuthorizeRequest,
		decodeGRPCAuthorizeResponse,
		pb.AuthorizeResponse{},
		grpctransport.ClientBefore(setMDCorrelationID),
	).Endpoint()

	validateEndpoint := grpctransport.NewClient(
		conn,
		"pb.Clients",
		"Validate",
		encodeGRPCValidateRequest,
		decodeGRPCValidateResponse,
		pb.ValidateResponse{},
		grpctransport.ClientBefore(setMDCorrelationID),
	).Endpoint()

	createEndpoint := grpctransport.NewClient(
		conn,
		"pb.Clients",
		"Create",
		encodeGRPCCreateRequest,
		decodeGRPCCreateResponse,
		pb.CreateResponse{},
		grpctransport.ClientBefore(setMDCorrelationID),
	).Endpoint()

	return Endpoints{
		AuthorizeEndpoint: authorizeEndpoint,
		ValidateEndpoint:  validateEndpoint,
		CreateEndpoint:    createEndpoint,
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

func encodeGRPCAuthorizeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding authorize request")
	req, ok := request.(*AuthorizeRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding: invalid authorize request structure")
		return &pb.AuthorizeRequest{}, ErrInvalidRequestStructure()
	}

	return &pb.AuthorizeRequest{
		User:     req.User,
		Password: req.Password,
	}, nil
}

func decodeGRPCAuthorizeResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding authorize response")
	reply, ok := grpcReply.(*pb.AuthorizeResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding: invalid authorize response structure")
		return &AuthorizeResponse{}, ErrInvalidResponseStructure()
	}
	return &AuthorizeResponse{
		Token: reply.Token,
		Err:   str2err(reply.Error),
	}, nil
}

func encodeGRPCValidateRequest(ctx context.Context, request interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding validate request")
	req, ok := request.(*ValidateRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding: invalid validate request structure")
		return &pb.ValidateRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.ValidateRequest{
		Token: req.Token,
	}, nil
}

func decodeGRPCValidateResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding validate response")
	reply, ok := grpcReply.(*pb.ValidateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding: invalid validate response structure")
		return &ValidateResponse{}, ErrInvalidResponseStructure()
	}
	return &ValidateResponse{
		User: reply.User,
		Err:  str2err(reply.Error),
	}, nil
}

func encodeGRPCCreateRequest(ctx context.Context, request interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding create request")
	req, ok := request.(*CreateRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding: invalid create request structure")
		return &pb.CreateRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.CreateRequest{
		User:     req.User,
		Password: req.Password,
	}, nil
}

func decodeGRPCCreateResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding create response")
	reply, ok := grpcReply.(*pb.CreateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding: invalid create response structure")
		return &CreateResponse{}, ErrInvalidResponseStructure()
	}
	return &CreateResponse{
		Err: str2err(reply.Error),
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
	case InvalidToken:
		return ErrInvalidToken()
	case UserNotFound:
		return ErrUserNotFound()
	case UserExists:
		return ErrUserExists()
	default:
		return errors.New(s)
	}
}
