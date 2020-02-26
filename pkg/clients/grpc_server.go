package clients

import (
	"context"
	"go-booking-service/commons"
	"go-booking-service/pb"

	"github.com/go-kit/kit/log/level"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type GrpcServer struct {
	authorize grpctransport.Handler
	validate  grpctransport.Handler
	create    grpctransport.Handler
}

func NewGRPCServer(endpoints Endpoints) *GrpcServer {
	return &GrpcServer{
		authorize: grpctransport.NewServer(
			endpoints.AuthorizeEndpoint,
			decodeGRPCAuthorizeRequest,
			encodeGRPCAuthorizeResponse,
			grpctransport.ServerBefore(setContextCorrelationId),
		),
		validate: grpctransport.NewServer(
			endpoints.ValidateEndpoint,
			decodeGRPCValidateRequest,
			encodeGRPCValidateResponse,
			grpctransport.ServerBefore(setContextCorrelationId),
		),
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
			grpctransport.ServerBefore(setContextCorrelationId),
		),
	}
}

func setContextCorrelationId(ctx context.Context, md metadata.MD) context.Context {
	if s, ok := md[commons.MetaDataKeyCorrelationID]; !ok || len(s) == 0 {
		level.Error(logger).Log("message", "setting metadata", "error", "no correlation id in metadata")
		return ctx
	}

	correlationID, err := uuid.Parse(md[commons.MetaDataKeyCorrelationID][0])
	if err != nil {
		level.Error(logger).Log("message", "setting metadata", "error", "invalid correlation id in metadata")
		return ctx
	}

	ctx = context.WithValue(ctx, commons.ContextKeyCorrelationID, correlationID)
	return ctx
}

func (s *GrpcServer) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "authorize request received")
	_, resp, err := s.authorize.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "authorize failed", "error", err)
		return &pb.AuthorizeResponse{}, err
	}
	response, ok := resp.(*pb.AuthorizeResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return &pb.AuthorizeResponse{}, ErrInvalidResponseStructure()
	}
	return response, nil
}

func (s *GrpcServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "validate request received")
	_, resp, err := s.validate.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "validate failed", "error", err)
		return &pb.ValidateResponse{}, err
	}
	response, ok := resp.(*pb.ValidateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return &pb.ValidateResponse{}, ErrInvalidResponseStructure()
	}
	return response, nil
}

func (s *GrpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "create request received")
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "check failed", "error", err)
		return &pb.CreateResponse{}, err
	}
	response, ok := resp.(*pb.CreateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return &pb.CreateResponse{}, ErrInvalidResponseStructure()
	}
	return response, nil
}

func decodeGRPCAuthorizeRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding authorize request")
	req, ok := grpcReq.(*pb.AuthorizeRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
		return &AuthorizeRequest{}, ErrInvalidRequestStructure()
	}
	return &AuthorizeRequest{
		User:     req.User,
		Password: req.Password,
	}, nil
}

func encodeGRPCAuthorizeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding authorize response")
	resp, ok := response.(*AuthorizeResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return &pb.AuthorizeResponse{}, ErrInvalidResponseStructure()
	}
	return &pb.AuthorizeResponse{
		Token: resp.Token,
		Error: err2str(resp.Err),
	}, nil
}

func decodeGRPCValidateRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding validate request")
	req, ok := grpcReq.(*pb.ValidateRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
		return &ValidateRequest{}, ErrInvalidRequestStructure()
	}
	return &ValidateRequest{
		Token: req.Token,
	}, nil
}

func encodeGRPCValidateResponse(ctx context.Context, response interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding validate response")
	resp, ok := response.(*ValidateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return &pb.ValidateResponse{}, ErrInvalidResponseStructure()
	}
	return &pb.ValidateResponse{
		User:  resp.User,
		Error: err2str(resp.Err),
	}, nil
}

func decodeGRPCCreateRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding create request")
	req, ok := grpcReq.(*pb.CreateRequest)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid request structure")
		return &CreateRequest{}, ErrInvalidRequestStructure()
	}
	return &CreateRequest{
		User:     req.User,
		Password: req.Password,
	}, nil
}

func encodeGRPCCreateResponse(ctx context.Context, response interface{}) (interface{}, error) {
	level.Debug(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding create response")
	resp, ok := response.(*CreateResponse)
	if !ok {
		level.Error(logger).Log("CorrelationID", ctx.Value(commons.ContextKeyCorrelationID), "message", "invalid response structure")
		return &pb.CreateResponse{}, ErrInvalidResponseStructure()
	}
	return &pb.CreateResponse{
		Error: err2str(resp.Err),
	}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
