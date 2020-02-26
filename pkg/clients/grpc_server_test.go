package clients

import (
	"context"
	"fmt"
	"go-booking-service/pb"
	"testing"

	"gotest.tools/assert"
)

type grpcAuthorizeCorrectMock struct{}

func (a grpcAuthorizeCorrectMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.AuthorizeResponse{Token: "jjj.www.ttt", Error: ""}, nil
}

type grpcValidateCorrectMock struct{}

func (a grpcValidateCorrectMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.ValidateResponse{User: "John", Error: ""}, nil
}

type grpcCreateCorrectMock struct{}

func (a grpcCreateCorrectMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.CreateResponse{Error: "error!"}, nil
}

type grpcAuthorizeWrongResponseMock struct{}

func (a grpcAuthorizeWrongResponseMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, "jjj.www.ttt", nil
}

type grpcValidateWrongResponseMock struct{}

func (a grpcValidateWrongResponseMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, "John", nil
}

type grpcCreateWrongResponseMock struct{}

func (a grpcCreateWrongResponseMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, fmt.Errorf("error!"), nil
}

type grpcAuthorizeErrorMock struct{}

func (a grpcAuthorizeErrorMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.AuthorizeResponse{}, ErrInvalidCredentials()
}

type grpcValidateErrorMock struct{}

func (a grpcValidateErrorMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.ValidateResponse{}, ErrInvalidToken()
}

type grpcCreateErrorMock struct{}

func (a grpcCreateErrorMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.CreateResponse{}, ErrInvalidCredentials()
}

var grpcServerAuthorizeTest = []struct {
	name    string
	server  *GrpcServer
	request *pb.AuthorizeRequest
	want    *pb.AuthorizeResponse
	err     error
}{
	{
		name: "should return the token",
		server: &GrpcServer{
			authorize: grpcAuthorizeCorrectMock{},
			validate:  grpcValidateCorrectMock{},
		},
		request: &pb.AuthorizeRequest{},
		want:    &pb.AuthorizeResponse{Token: "jjj.www.ttt", Error: ""},
	},
	{
		name: "should return en error if ServeGRPC returns an error",
		server: &GrpcServer{
			authorize: grpcAuthorizeErrorMock{},
			validate:  grpcValidateErrorMock{},
		},
		request: &pb.AuthorizeRequest{},
		want:    &pb.AuthorizeResponse{},
		err:     ErrInvalidCredentials(),
	},
	{
		name: "should return en error if the ServeGRPC response is the wrong type",
		server: &GrpcServer{
			authorize: grpcAuthorizeWrongResponseMock{},
			validate:  grpcValidateWrongResponseMock{},
		},
		request: &pb.AuthorizeRequest{},
		want:    &pb.AuthorizeResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestGRPCServerAuthorize(t *testing.T) {
	t.Log("GRPCServerAuthorize")

	for _, testcase := range grpcServerAuthorizeTest {
		t.Logf(testcase.name)

		result, err := testcase.server.Authorize(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var grpcServerValidateTest = []struct {
	name    string
	server  *GrpcServer
	request *pb.ValidateRequest
	want    *pb.ValidateResponse
	err     error
}{
	{
		name: "should return the token",
		server: &GrpcServer{
			authorize: grpcValidateCorrectMock{},
			validate:  grpcValidateCorrectMock{},
		},
		request: &pb.ValidateRequest{},
		want:    &pb.ValidateResponse{User: "John", Error: ""},
	},
	{
		name: "should return en error if ServeGRPC returns an error",
		server: &GrpcServer{
			authorize: grpcValidateErrorMock{},
			validate:  grpcValidateErrorMock{},
		},
		request: &pb.ValidateRequest{},
		want:    &pb.ValidateResponse{},
		err:     ErrInvalidToken(),
	},
	{
		name: "should return en error if the ServeGRPC response is the wrong type",
		server: &GrpcServer{
			authorize: grpcValidateWrongResponseMock{},
			validate:  grpcValidateWrongResponseMock{},
		},
		request: &pb.ValidateRequest{},
		want:    &pb.ValidateResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

var grpcServerCreateTest = []struct {
	name    string
	server  *GrpcServer
	request *pb.CreateRequest
	want    *pb.CreateResponse
	err     error
}{
	{
		name: "should return the error",
		server: &GrpcServer{
			authorize: grpcCreateCorrectMock{},
			validate:  grpcValidateCorrectMock{},
			create:    grpcCreateCorrectMock{},
		},
		request: &pb.CreateRequest{},
		want:    &pb.CreateResponse{Error: "error!"},
	},
	{
		name: "should return en error if ServeGRPC returns an error",
		server: &GrpcServer{
			authorize: grpcCreateErrorMock{},
			validate:  grpcValidateErrorMock{},
			create:    grpcCreateErrorMock{},
		},
		request: &pb.CreateRequest{},
		want:    &pb.CreateResponse{},
		err:     ErrInvalidCredentials(),
	},
	{
		name: "should return en error if the ServeGRPC response is the wrong type",
		server: &GrpcServer{
			authorize: grpcCreateWrongResponseMock{},
			validate:  grpcValidateWrongResponseMock{},
			create:    grpcCreateWrongResponseMock{},
		},
		request: &pb.CreateRequest{},
		want:    &pb.CreateResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestGRPCServerCreate(t *testing.T) {
	t.Log("GRPCServerCreate")

	for _, testcase := range grpcServerCreateTest {
		t.Logf(testcase.name)

		result, err := testcase.server.Create(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

func TestGRPCServerValidate(t *testing.T) {
	t.Log("GRPCServerValidate")

	for _, testcase := range grpcServerValidateTest {
		t.Logf(testcase.name)

		result, err := testcase.server.Validate(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var decodeGRPCAuthorizeRequestTest = []struct {
	name    string
	request interface{}
	want    *AuthorizeRequest
	err     error
}{
	{
		name:    "should return the new structure with the user information",
		request: &pb.AuthorizeRequest{User: "John", Password: "pass"},
		want:    &AuthorizeRequest{User: "John", Password: "pass"},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "jjj.www.ttt",
		want:    &AuthorizeRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestDecodeGRPCAuthorizeRequest(t *testing.T) {
	t.Log("decodeGRPCAuthorizeRequest")

	for _, testcase := range decodeGRPCAuthorizeRequestTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCAuthorizeRequest(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var decodeGRPCValidateRequestTest = []struct {
	name    string
	request interface{}
	want    *ValidateRequest
	err     error
}{
	{
		name:    "should return the new structure with the user",
		request: &pb.ValidateRequest{Token: "Jhon"},
		want:    &ValidateRequest{Token: "Jhon"},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "Jhon",
		want:    &ValidateRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestDecodeGRPCValidateRequest(t *testing.T) {
	t.Log("decodeGRPCValidateRequest")

	for _, testcase := range decodeGRPCValidateRequestTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCValidateRequest(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var decodeGRPCCreateRequestTest = []struct {
	name    string
	request interface{}
	want    *CreateRequest
	err     error
}{
	{
		name:    "should return the new structure with the user and password",
		request: &pb.CreateRequest{User: "John", Password: "pass"},
		want:    &CreateRequest{User: "John", Password: "pass"},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "John, pass",
		want:    &CreateRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestDecodeGRPCCreateRequest(t *testing.T) {
	t.Log("decodeGRPCCreateRequest")

	for _, testcase := range decodeGRPCCreateRequestTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCCreateRequest(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var encodeGRPCAuthorizeResponseTest = []struct {
	name    string
	request interface{}
	want    *pb.AuthorizeResponse
	err     error
}{
	{
		name:    "should return the user",
		request: &AuthorizeResponse{Token: "jjj.www.ttt", Err: nil},
		want:    &pb.AuthorizeResponse{Token: "jjj.www.ttt", Error: ""},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "jjj.www.ttt",
		want:    &pb.AuthorizeResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestEncodeGRPCAuthorizeResponse(t *testing.T) {
	t.Log("EncodeGRPCAuthorizeResponse")

	for _, testcase := range encodeGRPCAuthorizeResponseTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCAuthorizeResponse(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var encodeGRPCValidateResponseTest = []struct {
	name    string
	request interface{}
	want    *pb.ValidateResponse
	err     error
}{
	{
		name:    "should return the user",
		request: &ValidateResponse{User: "John"},
		want:    &pb.ValidateResponse{User: "John"},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "John",
		want:    &pb.ValidateResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestEncodeGRPCValidateResponse(t *testing.T) {
	t.Log("encodeGRPCValidateResponse")

	for _, testcase := range encodeGRPCValidateResponseTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCValidateResponse(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var encodeGRPCCreateResponseTest = []struct {
	name    string
	request interface{}
	want    *pb.CreateResponse
	err     error
}{
	{
		name:    "should return the error",
		request: &CreateResponse{Err: fmt.Errorf("error!")},
		want:    &pb.CreateResponse{Error: "error!"},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "error!",
		want:    &pb.CreateResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestEncodeGRPCCreateResponse(t *testing.T) {
	t.Log("EncodeGRPCCreateResponse")

	for _, testcase := range encodeGRPCCreateResponseTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCCreateResponse(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}
