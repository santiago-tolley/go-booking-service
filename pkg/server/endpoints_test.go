package server

import (
	"context"
	"go-booking-service/pkg/clients"
	"go-booking-service/pkg/rooms"
	"testing"
	"time"

	"github.com/go-kit/kit/endpoint"
	"gotest.tools/assert"
)

type mockCorrectEndpoint struct{}

func (m mockCorrectEndpoint) Authorize(ctx context.Context, user, password string) (string, error) {
	return "jjj.www.ttt", nil
}

func (m mockCorrectEndpoint) Validate(ctx context.Context, token string) (string, error) {
	return "Jhon", nil
}

func (m mockCorrectEndpoint) Book(ctx context.Context, token string, date time.Time) (int, error) {
	return 1, nil
}

func (m mockCorrectEndpoint) Check(ctx context.Context, date time.Time) (int, error) {
	return 5, nil
}

type mockErrorEndpoint struct{}

func (m mockErrorEndpoint) Authorize(ctx context.Context, user, password string) (string, error) {
	return "", clients.ErrInvalidCredentials()
}

func (m mockErrorEndpoint) Validate(ctx context.Context, token string) (string, error) {
	return "", clients.ErrUserNotFound()
}

func (m mockErrorEndpoint) Book(ctx context.Context, token string, date time.Time) (int, error) {
	return 0, rooms.ErrNoRoomAvailable()
}

func (m mockErrorEndpoint) Check(ctx context.Context, date time.Time) (int, error) {
	return 0, rooms.ErrNoRoomAvailable()
}

type mockInvalidEndpoint struct{}

func (m mockInvalidEndpoint) Authorize(ctx context.Context, user, password string) (string, error) {
	return "", ErrInvalidResponseStructure()
}

func (m mockInvalidEndpoint) Validate(ctx context.Context, token string) (string, error) {
	return "", ErrInvalidResponseStructure()
}

func (m mockInvalidEndpoint) Book(ctx context.Context, token string, date time.Time) (int, error) {
	return 0, rooms.ErrInvalidResponseStructure()
}

func (m mockInvalidEndpoint) Check(ctx context.Context, date time.Time) (int, error) {
	return 0, rooms.ErrInvalidResponseStructure()
}

var endpointAuthorizeTest = []struct {
	name              string
	user              string
	password          string
	authorizeEndpoint endpoint.Endpoint
	want              string
	err               error
}{
	{
		name:     "should return the token",
		user:     "Jhon",
		password: "pass",
		authorizeEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &AuthorizeResponse{"jjj.www.ttt", nil}, nil
		},
		want: "jjj.www.ttt",
	},
	{
		name:     "should return an error if the response has the wrong structure",
		user:     "Charles",
		password: "pass",
		authorizeEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return "jjj.www.ttt", nil
		},
		want: "",
		err:  ErrInvalidResponseStructure(),
	},
	{
		name:     "should return an error if the endpoint returns an error",
		user:     "Jhon",
		password: "pass",
		authorizeEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &AuthorizeResponse{}, clients.ErrInvalidCredentials()
		},
		want: "",
		err:  clients.ErrInvalidCredentials(),
	},
}

func TestEndpointAuthorize(t *testing.T) {
	t.Log("EndpointAuthorize")

	for _, testcase := range endpointAuthorizeTest {
		t.Logf(testcase.name)

		endpointMock := Endpoints{
			AuthorizeEndpoint: testcase.authorizeEndpoint,
		}
		result, err := endpointMock.Authorize(context.Background(), testcase.user, testcase.password)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var endpointValidateTest = []struct {
	name             string
	token            string
	validateEndpoint endpoint.Endpoint
	want             string
	err              error
}{
	{
		name:  "should return the user",
		token: "jjj.www.ttt",
		validateEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &ValidateResponse{"John", nil}, nil
		},
		want: "John",
	},
	{
		name:  "should return an error if the response has the wrong structure",
		token: "jjj.www.ttt",
		validateEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return "John", nil
		},
		want: "",
		err:  ErrInvalidResponseStructure(),
	},
	{
		name:  "should return an error if the endpoint returns an error",
		token: "jjj.www.ttt",
		validateEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return nil, clients.ErrInvalidToken()
		},
		want: "",
		err:  clients.ErrInvalidToken(),
	},
}

func TestEndpointValidate(t *testing.T) {
	t.Log("EndpointValidate")

	for _, testcase := range endpointValidateTest {
		t.Logf(testcase.name)

		endpointMock := Endpoints{
			ValidateEndpoint: testcase.validateEndpoint,
		}
		result, err := endpointMock.Validate(context.Background(), testcase.token)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var endpointBookTest = []struct {
	name         string
	token        string
	date         time.Time
	bookEndpoint endpoint.Endpoint
	want         int
	err          error
}{
	{
		name:  "should return the room id",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		bookEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &BookResponse{1, nil}, nil
		},
		want: 1,
	},
	{
		name:  "should return an error if the response has the wrong structure",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		bookEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return 1, nil
		},
		want: 0,
		err:  ErrInvalidResponseStructure(),
	},
	{
		name:  "should return an error if the endpoint returns an error",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		bookEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return nil, rooms.ErrNoRoomAvailable()
		},
		want: 0,
		err:  rooms.ErrNoRoomAvailable(),
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

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
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
		name: "should return the number of available rooms",
		date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		checkEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &CheckResponse{5, nil}, nil
		},
		want: 5,
	},
	{
		name: "should return an error if the response has the wrong structure",
		date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		checkEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return 5, nil
		},
		want: 0,
		err:  ErrInvalidResponseStructure(),
	},
	{
		name: "should return an error if the endpoint returns an error",
		date: time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		checkEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return nil, rooms.ErrNoRoomAvailable()
		},
		want: 0,
		err:  rooms.ErrNoRoomAvailable(),
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

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}
