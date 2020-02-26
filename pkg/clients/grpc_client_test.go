package clients

import (
	"context"
	"go-booking-service/pb"
	"testing"

	"gotest.tools/assert"
)

var encodeGRPCAuthorizeRequestTest = []struct {
	name    string
	request interface{}
	want    *pb.AuthorizeRequest
	err     error
}{
	{
		name:    "should return the user",
		request: &AuthorizeRequest{User: "John", Password: "pass"},
		want:    &pb.AuthorizeRequest{User: "John", Password: "pass"},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "John",
		want:    &pb.AuthorizeRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestEncodeGRPCAuthorizeRequest(t *testing.T) {
	t.Log("EncodeGRPCAuthorizeRequest")

	for _, testcase := range encodeGRPCAuthorizeRequestTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCAuthorizeRequest(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var encodeGRPCValidateRequestTest = []struct {
	name    string
	request interface{}
	want    *pb.ValidateRequest
	err     error
}{
	{
		name:    "should return the user",
		request: &ValidateRequest{Token: "jjj.www.ttt"},
		want:    &pb.ValidateRequest{Token: "jjj.www.ttt"},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "jjj.www.ttt",
		want:    &pb.ValidateRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestEncodeGRPCValidateRequest(t *testing.T) {
	t.Log("encodeGRPCValidateRequest")

	for _, testcase := range encodeGRPCValidateRequestTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCValidateRequest(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var decodeGRPCAuthorizeResponseTest = []struct {
	name    string
	request interface{}
	want    *AuthorizeResponse
	err     error
}{
	{
		name:    "should return the new structure with the token",
		request: &pb.AuthorizeResponse{Token: "jjj.www.ttt", Error: ""},
		want:    &AuthorizeResponse{Token: "jjj.www.ttt", Err: nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "jjj.www.ttt",
		want:    &AuthorizeResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestDecodeGRPCAuthorizeResponse(t *testing.T) {
	t.Log("decodeGRPCAuthorizeResponse")

	for _, testcase := range decodeGRPCAuthorizeResponseTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCAuthorizeResponse(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var decodeGRPCValidateResponseTest = []struct {
	name    string
	request interface{}
	want    *ValidateResponse
	err     error
}{
	{
		name:    "should return the new structure with the user",
		request: &pb.ValidateResponse{User: "John", Error: ""},
		want:    &ValidateResponse{User: "John", Err: nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "John",
		want:    &ValidateResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestDecodeGRPCValidateResponse(t *testing.T) {
	t.Log("decodeGRPCValidateResponse")

	for _, testcase := range decodeGRPCValidateResponseTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCValidateResponse(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}
