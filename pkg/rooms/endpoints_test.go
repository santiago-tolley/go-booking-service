package rooms

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-kit/kit/endpoint"
)

var endpointBookTest = []struct {
	name         string
	token        string
	date         time.Time
	bookEndpoint endpoint.Endpoint
	want         int
	err          error
}{
	{
		name:  "should return booked room id",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		bookEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &BookResponse{5, nil}, nil
		},
		want: 5,
	},
	{
		name:  "should return an error if the endpoint returns an error",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		bookEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &BookResponse{}, ErrNoRoomAvailable()
		},
		err: ErrNoRoomAvailable(),
	},
	{
		name:  "should return an error if response structure is incorrect",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		bookEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return 5, nil
		},
		err: ErrInvalidResponseStructure(),
	},
}

func TestEndpointBook(t *testing.T) {
	t.Log("EndpointBook")

	for _, testcase := range endpointBookTest {
		t.Logf(testcase.name)

		endpointMock := Endpoints{
			BookEndpoint: testcase.bookEndpoint,
		}
		result, err := endpointMock.Book(context.Background(), testcase.token, testcase.date)

		if result != testcase.want {
			t.Errorf("=> Got %v wanted %v", result, testcase.want)
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

var endpointCheckTest = []struct {
	name          string
	date          time.Time
	checkEndpoint endpoint.Endpoint
	want          int
	err           error
}{
	{
		name: "should return number of available rooms",
		date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		checkEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &CheckResponse{5, nil}, nil
		},
		want: 5,
	},
	{
		name: "should return an error if the endpoint returns an error",
		date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		checkEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &CheckResponse{}, ErrNoRoomAvailable()
		},
		err: ErrNoRoomAvailable(),
	},
	{
		name: "should return an error if response structure is incorrect",
		date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		checkEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return 5, nil
		},
		err: ErrInvalidResponseStructure(),
	},
}

func TestEndpointCheck(t *testing.T) {
	t.Log("EndpointCheck")

	for _, testcase := range endpointCheckTest {
		t.Logf(testcase.name)

		endpointMock := Endpoints{
			CheckEndpoint: testcase.checkEndpoint,
		}
		result, err := endpointMock.Check(context.Background(), testcase.date)

		if result != testcase.want {
			t.Errorf("=> Got %v wanted %v", result, testcase.want)
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

type mockCorrectClientsService struct{}

func (m mockCorrectClientsService) Book(ctx context.Context, token string, daet time.Time) (int, error) {
	return 1, nil
}

func (m mockCorrectClientsService) Check(ctx context.Context, date time.Time) (int, error) {
	return 5, nil
}

type mockErrorClientsService struct{}

func (m mockErrorClientsService) Book(ctx context.Context, token string, daet time.Time) (int, error) {
	return 0, ErrNoRoomAvailable()
}

func (m mockErrorClientsService) Check(ctx context.Context, date time.Time) (int, error) {
	return 0, ErrNoRoomAvailable()
}

var makeBookEndpointTest = []struct {
	name    string
	client  RoomsService
	request interface{}
	want    *BookResponse
	err     error
}{
	{
		name:    "should return the booked room id",
		client:  mockCorrectClientsService{},
		request: &BookRequest{},
		want:    &BookResponse{1, nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		client:  mockCorrectClientsService{},
		request: "jjj.www.ttt",
		want:    &BookResponse{},
		err:     ErrInvalidRequestStructure(),
	},
	{
		name:    "should return an error if the endpoint returns an error",
		client:  mockErrorClientsService{},
		request: &BookRequest{},
		want:    &BookResponse{0, ErrNoRoomAvailable()},
	},
}

func TestMakeBookEndpoint(t *testing.T) {
	t.Log("MakeBookEndpoint")

	for _, testcase := range makeBookEndpointTest {
		t.Logf(testcase.name)

		endpoint := MakeBookEndpoint(testcase.client)
		result, err := endpoint(context.Background(), testcase.request)

		if !reflect.DeepEqual(result.(*BookResponse), testcase.want) {
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

var makeCheckEndpointTest = []struct {
	name    string
	client  RoomsService
	request interface{}
	want    *CheckResponse
	err     error
}{
	{
		name:    "should return the booked room id",
		client:  mockCorrectClientsService{},
		request: &CheckRequest{},
		want:    &CheckResponse{5, nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		client:  mockCorrectClientsService{},
		request: "jjj.www.ttt",
		want:    &CheckResponse{},
		err:     ErrInvalidRequestStructure(),
	},
	{
		name:    "should return an error if the endpoint returns an error",
		client:  mockErrorClientsService{},
		request: &CheckRequest{},
		want:    &CheckResponse{0, ErrNoRoomAvailable()},
	},
}

func TestMakeCheckEndpoint(t *testing.T) {
	t.Log("MakeCheckEndpoint")

	for _, testcase := range makeCheckEndpointTest {
		t.Logf(testcase.name)

		endpoint := MakeCheckEndpoint(testcase.client)
		result, err := endpoint(context.Background(), testcase.request)

		if !reflect.DeepEqual(result.(*CheckResponse), testcase.want) {
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
