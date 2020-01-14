package clients

import (
	"context"
	"go-booking-service/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type GrpcServer struct {
	authorize grpctransport.Handler
	validate  grpctransport.Handler
}

func NewGRPCServer(endpoints Endpoints) *GrpcServer {
	return &GrpcServer{
		authorize: grpctransport.NewServer(
			endpoints.AuthorizeEndpoint,
			decodeGRPCAuthorizeRequest,
			encodeGRPCAuthorizeResponse,
		),
		validate: grpctransport.NewServer(
			endpoints.ValidateEndpoint,
			decodeGRPCValidateRequest,
			encodeGRPCValidateResponse,
		),
	}
}

func (s *GrpcServer) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	_, resp, err := s.authorize.ServeGRPC(ctx, req)
	if err != nil {
		return &pb.AuthorizeResponse{}, err
	}
	response, ok := resp.(*pb.AuthorizeResponse)
	if !ok {
		return &pb.AuthorizeResponse{}, ErrInvalidResponseStructure()
	}
	return response, nil
}

func (s *GrpcServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, req)
	if err != nil {
		return &pb.ValidateResponse{}, err
	}
	response, ok := resp.(*pb.ValidateResponse)
	if !ok {
		return &pb.ValidateResponse{}, ErrInvalidResponseStructure()
	}
	return response, nil
}

func decodeGRPCAuthorizeRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.AuthorizeRequest)
	if !ok {
		return AuthorizeRequest{}, ErrInvalidRequestStructure()
	}
	return AuthorizeRequest{
		User:     req.User,
		Password: req.Password,
	}, nil
}

func encodeGRPCAuthorizeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(AuthorizeResponse)
	if !ok {
		return &pb.AuthorizeResponse{}, ErrInvalidResponseStructure()
	}
	return &pb.AuthorizeResponse{
		Token: resp.Token,
		Error: err2str(resp.Err),
	}, nil
}

func decodeGRPCValidateRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.ValidateRequest)
	if !ok {
		return ValidateRequest{}, ErrInvalidRequestStructure()
	}
	return ValidateRequest{
		Token: req.Token,
	}, nil
}

func encodeGRPCValidateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(ValidateResponse)
	if !ok {
		return &pb.ValidateResponse{}, ErrInvalidResponseStructure()
	}
	return &pb.ValidateResponse{
		User:  resp.User,
		Error: err2str(resp.Err),
	}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
