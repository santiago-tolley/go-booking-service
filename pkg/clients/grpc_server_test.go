package clients

import (
	"context"
	"go-booking-service/pb"
	"reflect"
	"testing"
)

type grpcAuthorizeCorrectMock struct{}

func (a grpcAuthorizeCorrectMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.AuthorizeResponse{Token: "jjj.www.ttt", Error: ""}, nil
}

type grpcValidateCorrectMock struct{}

func (a grpcValidateCorrectMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.ValidateResponse{User: "John", Error: ""}, nil
}

type grpcAuthorizeWrongResponseMock struct{}

func (a grpcAuthorizeWrongResponseMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, "jjj.www.ttt", nil
}

type grpcValidateWrongResponseMock struct{}

func (a grpcValidateWrongResponseMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, "John", nil
}

type grpcAuthorizeErrorMock struct{}

func (a grpcAuthorizeErrorMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.AuthorizeResponse{}, ErrInvalidCredentials()
}

type grpcValidateErrorMock struct{}

func (a grpcValidateErrorMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.ValidateResponse{}, ErrInvalidToken()
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

		if !reflect.DeepEqual(result, testcase.want) {
			t.Errorf("=> Got %v (%T) wanted %v (%T)", result, result, testcase.want, testcase.want)
		}

		var ok bool
		if testcase.err != nil {
			if err == testcase.err {
				ok = true
			}
		} else if err == nil {
			ok = true
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
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

func TestGRPCServerValidate(t *testing.T) {
	t.Log("GRPCServerValidate")

	for _, testcase := range grpcServerValidateTest {
		t.Logf(testcase.name)

		result, err := testcase.server.Validate(context.Background(), testcase.request)

		if !reflect.DeepEqual(result, testcase.want) {
			t.Errorf("=> Got %v (%T) wanted %v (%T)", result, result, testcase.want, testcase.want)
		}

		var ok bool
		if testcase.err != nil {
			if err == testcase.err {
				ok = true
			}
		} else if err == nil {
			ok = true
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
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

		if !reflect.DeepEqual(result.(*AuthorizeRequest), testcase.want) {
			t.Errorf("=> Got %v (%T) wanted %v (%T)", result, result, testcase.want, testcase.want)
		}

		var ok bool
		if testcase.err != nil {
			if err == testcase.err {
				ok = true
			}
		} else if err == nil {
			ok = true
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
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

		if !reflect.DeepEqual(result.(*ValidateRequest), testcase.want) {
			t.Errorf("=> Got %v (%T) wanted %v (%T)", result, result, testcase.want, testcase.want)
		}

		var ok bool
		if testcase.err != nil {
			if err == testcase.err {
				ok = true
			}
		} else if err == nil {
			ok = true
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
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

		if !reflect.DeepEqual(result.(*pb.AuthorizeResponse), testcase.want) {
			t.Errorf("=> Got %v (%T) wanted %v (%T)", result, result, testcase.want, testcase.want)
		}

		var ok bool
		if testcase.err != nil {
			if err == testcase.err {
				ok = true
			}
		} else if err == nil {
			ok = true
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
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

		if !reflect.DeepEqual(result.(*pb.ValidateResponse), testcase.want) {
			t.Errorf("=> Got %v (%T) wanted %v (%T)", result, result, testcase.want, testcase.want)
		}

		var ok bool
		if testcase.err != nil {
			if err == testcase.err {
				ok = true
			}
		} else if err == nil {
			ok = true
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
	}
}
