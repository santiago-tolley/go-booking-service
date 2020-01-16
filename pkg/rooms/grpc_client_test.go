package rooms

import (
	"context"
	"go-booking-service/pb"
	"testing"
	"time"

	"gotest.tools/assert"
)

var encodeGRPCBookRequestTest = []struct {
	name    string
	request interface{}
	want    *pb.BookRequest
	err     error
}{
	{
		name:    "should return the values in the pb structure",
		request: &BookRequest{Token: "jjj.www.ttt", Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC)},
		want:    &pb.BookRequest{Token: "jjj.www.ttt", Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC).Unix()},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "jjj.www.ttt",
		want:    &pb.BookRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestEncodeGRPCBookRequest(t *testing.T) {
	t.Log("EncodeGRPCBookRequest")

	for _, testcase := range encodeGRPCBookRequestTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCBookRequest(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var encodeGRPCCheckRequestTest = []struct {
	name    string
	request interface{}
	want    *pb.CheckRequest
	err     error
}{
	{
		name:    "should return the values in the pb structure",
		request: &CheckRequest{Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC)},
		want:    &pb.CheckRequest{Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC).Unix()},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		want:    &pb.CheckRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestEncodeGRPCCheckRequest(t *testing.T) {
	t.Log("encodeGRPCCheckRequest")

	for _, testcase := range encodeGRPCCheckRequestTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCCheckRequest(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var decodeGRPCBookResponseTest = []struct {
	name    string
	request interface{}
	want    *BookResponse
	err     error
}{
	{
		name:    "should return the values in the internal structure",
		request: &pb.BookResponse{Id: 1, Error: ""},
		want:    &BookResponse{Id: 1, Err: nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "jjj.www.ttt",
		want:    &BookResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestDecodeGRPCBookResponse(t *testing.T) {
	t.Log("decodeGRPCBookResponse")

	for _, testcase := range decodeGRPCBookResponseTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCBookResponse(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var decodeGRPCCheckResponseTest = []struct {
	name    string
	request interface{}
	want    *CheckResponse
	err     error
}{
	{
		name:    "should return the values in the internal structure",
		request: &pb.CheckResponse{Available: 5, Error: ""},
		want:    &CheckResponse{Available: 5, Err: nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "Jhon",
		want:    &CheckResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestDecodeGRPCCheckResponse(t *testing.T) {
	t.Log("decodeGRPCCheckResponse")

	for _, testcase := range decodeGRPCCheckResponseTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCCheckResponse(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}
