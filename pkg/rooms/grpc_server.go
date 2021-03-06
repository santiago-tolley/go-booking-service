package rooms

import (
	"context"
	"go-booking-service/pb"
	"time"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type GrpcServer struct {
	book  grpctransport.Handler
	check grpctransport.Handler
}

func NewGRPCServer(endpoints Endpoints) *GrpcServer {
	return &GrpcServer{
		book: grpctransport.NewServer(
			endpoints.BookEndpoint,
			decodeGRPCBookRequest,
			encodeGRPCBookResponse,
		),
		check: grpctransport.NewServer(
			endpoints.CheckEndpoint,
			decodeGRPCCheckRequest,
			encodeGRPCCheckResponse,
		),
	}
}

func (s *GrpcServer) Book(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	_, resp, err := s.book.ServeGRPC(ctx, req)
	if err != nil {
		return &pb.BookResponse{}, err
	}
	response, ok := resp.(*pb.BookResponse)
	if !ok {
		return &pb.BookResponse{}, ErrInvalidResponseStructure()
	}
	return response, nil
}

func (s *GrpcServer) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	_, resp, err := s.check.ServeGRPC(ctx, req)
	if err != nil {
		return &pb.CheckResponse{}, err
	}
	response, ok := resp.(*pb.CheckResponse)
	if !ok {
		return &pb.CheckResponse{}, ErrInvalidResponseStructure()
	}
	return response, nil
}

func decodeGRPCBookRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.BookRequest)
	if !ok {
		return &BookRequest{}, ErrInvalidRequestStructure()
	}
	return &BookRequest{
		Token: req.Token,
		Date:  time.Unix(req.Date, 0).UTC(),
	}, nil
}

func decodeGRPCCheckRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.CheckRequest)
	if !ok {
		return &CheckRequest{}, ErrInvalidRequestStructure()
	}
	return &CheckRequest{
		Date: time.Unix(req.Date, 0).UTC(),
	}, nil
}

func encodeGRPCBookResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*BookResponse)
	if !ok {
		return &pb.BookResponse{}, ErrInvalidResponseStructure()
	}
	return &pb.BookResponse{
		Id:    int64(resp.Id),
		Error: err2str(resp.Err),
	}, nil
}

func encodeGRPCCheckResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*CheckResponse)
	if !ok {
		return &pb.CheckResponse{}, ErrInvalidResponseStructure()
	}
	return &pb.CheckResponse{
		Available: int64(resp.Available),
		Error:     err2str(resp.Err),
	}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
