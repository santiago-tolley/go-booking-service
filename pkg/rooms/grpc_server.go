package rooms

import (
	"context"
	"go-booking-service/commons"
	"go-booking-service/pb"
	"time"

	"github.com/go-kit/kit/log/level"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
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
			grpctransport.ServerBefore(setContextCorrelationId),
		),
		check: grpctransport.NewServer(
			endpoints.CheckEndpoint,
			decodeGRPCCheckRequest,
			encodeGRPCCheckResponse,
			grpctransport.ServerBefore(setContextCorrelationId),
		),
	}
}

func setContextCorrelationId(ctx context.Context, md metadata.MD) context.Context {
	if s, ok := md[commons.MetaDataKeyCorrelationID]; !ok || len(s) == 0 {
		level.Error(logger).Log("message", "setting metadata", "error", "no correlation id in metadata")
		return ctx
	}

	level.Debug(logger).Log("message", "setting context", "metadata", md[commons.MetaDataKeyCorrelationID][0])
	correlationID, err := uuid.Parse(md[commons.MetaDataKeyCorrelationID][0])
	if err != nil {
		level.Error(logger).Log("message", "setting metadata", "error", "invalid correlation id in metadata")
		return ctx
	}
	ctx = context.WithValue(ctx, commons.ContextKeyCorrelationID, correlationID)
	return ctx
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
