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
	level.Debug(logger).Log("message", "setMDCorrelationID", "ctx", ctx.Value(commons.ContextKeyCorrelationID))
	correlationId, ok := ctx.Value(commons.ContextKeyCorrelationID).(uuid.UUID)
	if !ok {
		level.Error(logger).Log("message", "setting metadata", "error", "no correlation id in context")
		return ctx
	}

	level.Debug(logger).Log("message", "setting metadata", "context", correlationId.String())
	(*md)[commons.MetaDataKeyCorrelationID] = append((*md)[commons.MetaDataKeyCorrelationID], correlationId.String())
	return ctx
}

func encodeGRPCAuthorizeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*AuthorizeRequest)
	if !ok {
		return &pb.AuthorizeRequest{}, ErrInvalidRequestStructure()
	}

	return &pb.AuthorizeRequest{
		User:     req.User,
		Password: req.Password,
	}, nil
}

func decodeGRPCAuthorizeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.AuthorizeResponse)
	if !ok {
		return &AuthorizeResponse{}, ErrInvalidResponseStructure()
	}
	return &AuthorizeResponse{
		Token: reply.Token,
		Err:   str2err(reply.Error),
	}, nil
}

func encodeGRPCValidateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*ValidateRequest)
	if !ok {
		return &pb.ValidateRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.ValidateRequest{
		Token: req.Token,
	}, nil
}

func decodeGRPCValidateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.ValidateResponse)
	if !ok {
		return &ValidateResponse{}, ErrInvalidResponseStructure()
	}
	return &ValidateResponse{
		User: reply.User,
		Err:  str2err(reply.Error),
	}, nil
}

func encodeGRPCCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*CreateRequest)
	if !ok {
		return &pb.CreateRequest{}, ErrInvalidRequestStructure()
	}
	return &pb.CreateRequest{
		User:     req.User,
		Password: req.Password,
	}, nil
}

func decodeGRPCCreateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, ok := grpcReply.(*pb.CreateResponse)
	if !ok {
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
