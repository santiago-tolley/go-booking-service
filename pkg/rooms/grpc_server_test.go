package rooms

import (
	"context"
	"go-booking-service/pb"
	"reflect"
	"testing"
	"time"
)

type grpcBookCorrectMock struct{}

func (b grpcBookCorrectMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.BookResponse{Id: 1, Error: ""}, nil
}

type grpcCheckCorrectMock struct{}

func (a grpcCheckCorrectMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.CheckResponse{Available: 5, Error: ""}, nil
}

type grpcBookWrongResponseMock struct{}

func (a grpcBookWrongResponseMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, 1, nil
}

type grpcCheckWrongResponseMock struct{}

func (a grpcCheckWrongResponseMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, 5, nil
}

type grpcBookErrorMock struct{}

func (a grpcBookErrorMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.BookResponse{}, ErrNoRoomAvailable()
}

type grpcCheckErrorMock struct{}

func (a grpcCheckErrorMock) ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error) {
	return ctx, &pb.CheckResponse{}, ErrNoRoomAvailable()
}

var grpcServerBookTest = []struct {
	name    string
	server  *GrpcServer
	request *pb.BookRequest
	want    *pb.BookResponse
	err     error
}{
	{
		name: "should return the booked room id",
		server: &GrpcServer{
			book:  grpcBookCorrectMock{},
			check: grpcCheckCorrectMock{},
		},
		request: &pb.BookRequest{},
		want:    &pb.BookResponse{Id: 1, Error: ""},
	},
	{
		name: "should return en error if ServeGRPC returns an error",
		server: &GrpcServer{
			book:  grpcBookErrorMock{},
			check: grpcCheckErrorMock{},
		},
		request: &pb.BookRequest{},
		want:    &pb.BookResponse{},
		err:     ErrNoRoomAvailable(),
	},
	{
		name: "should return en error if the ServeGRPC response is the wrong type",
		server: &GrpcServer{
			book:  grpcBookWrongResponseMock{},
			check: grpcCheckWrongResponseMock{},
		},
		request: &pb.BookRequest{},
		want:    &pb.BookResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestGRPCServerBook(t *testing.T) {
	t.Log("GRPCServerBook")

	for _, testcase := range grpcServerBookTest {
		t.Logf(testcase.name)

		result, err := testcase.server.Book(context.Background(), testcase.request)

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

var grpcServerCheckTest = []struct {
	name    string
	server  *GrpcServer
	request *pb.CheckRequest
	want    *pb.CheckResponse
	err     error
}{
	{
		name: "should return the number of available rooms",
		server: &GrpcServer{
			book:  grpcBookCorrectMock{},
			check: grpcCheckCorrectMock{},
		},
		request: &pb.CheckRequest{},
		want:    &pb.CheckResponse{Available: 5, Error: ""},
	},
	{
		name: "should return en error if ServeGRPC returns an error",
		server: &GrpcServer{
			book:  grpcBookErrorMock{},
			check: grpcCheckErrorMock{},
		},
		request: &pb.CheckRequest{},
		want:    &pb.CheckResponse{},
		err:     ErrNoRoomAvailable(),
	},
	{
		name: "should return en error if the ServeGRPC response is the wrong type",
		server: &GrpcServer{
			book:  grpcBookWrongResponseMock{},
			check: grpcCheckWrongResponseMock{},
		},
		request: &pb.CheckRequest{},
		want:    &pb.CheckResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestGRPCServerCheck(t *testing.T) {
	t.Log("GRPCServerCheck")

	for _, testcase := range grpcServerCheckTest {
		t.Logf(testcase.name)

		result, err := testcase.server.Check(context.Background(), testcase.request)

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

var decodeGRPCBookRequestTest = []struct {
	name    string
	request interface{}
	want    BookRequest
	err     error
}{
	{
		name:    "should return values in the internal structure",
		request: &pb.BookRequest{Token: "jjj.www.ttt", Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC).Unix()},
		want:    BookRequest{Token: "jjj.www.ttt", Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC)},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "jjj.www.ttt",
		want:    BookRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestDecodeGRPCBookRequest(t *testing.T) {
	t.Log("decodeGRPCBookRequest")

	for _, testcase := range decodeGRPCBookRequestTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCBookRequest(context.Background(), testcase.request)

		if !reflect.DeepEqual(result.(BookRequest), testcase.want) {
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

var decodeGRPCCheckRequestTest = []struct {
	name    string
	request interface{}
	want    CheckRequest
	err     error
}{
	{
		name:    "should return the new structure with the date",
		request: &pb.CheckRequest{Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC).Unix()},
		want:    CheckRequest{Date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC)},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC).Unix(),
		want:    CheckRequest{},
		err:     ErrInvalidRequestStructure(),
	},
}

func TestDecodeGRPCCheckRequest(t *testing.T) {
	t.Log("decodeGRPCCheckRequest")

	for _, testcase := range decodeGRPCCheckRequestTest {
		t.Logf(testcase.name)

		result, err := decodeGRPCCheckRequest(context.Background(), testcase.request)

		if !reflect.DeepEqual(result.(CheckRequest), testcase.want) {
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

var encodeGRPCBookResponseTest = []struct {
	name    string
	request interface{}
	want    *pb.BookResponse
	err     error
}{
	{
		name:    "should return the pb structure with the id",
		request: BookResponse{Id: 1, Err: nil},
		want:    &pb.BookResponse{Id: 1, Error: ""},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: 1,
		want:    &pb.BookResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestEncodeGRPCBookResponse(t *testing.T) {
	t.Log("EncodeGRPCBookResponse")

	for _, testcase := range encodeGRPCBookResponseTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCBookResponse(context.Background(), testcase.request)

		if !reflect.DeepEqual(result.(*pb.BookResponse), testcase.want) {
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

var encodeGRPCCheckResponseTest = []struct {
	name    string
	request interface{}
	want    *pb.CheckResponse
	err     error
}{
	{
		name:    "should return the pb structure with the number of available rooms",
		request: CheckResponse{Available: 5},
		want:    &pb.CheckResponse{Available: 5},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		request: "John",
		want:    &pb.CheckResponse{},
		err:     ErrInvalidResponseStructure(),
	},
}

func TestEncodeGRPCCheckResponse(t *testing.T) {
	t.Log("encodeGRPCCheckResponse")

	for _, testcase := range encodeGRPCCheckResponseTest {
		t.Logf(testcase.name)

		result, err := encodeGRPCCheckResponse(context.Background(), testcase.request)

		if !reflect.DeepEqual(result.(*pb.CheckResponse), testcase.want) {
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
