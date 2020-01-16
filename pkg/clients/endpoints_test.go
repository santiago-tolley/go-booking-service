package clients

import (
	"context"
	"fmt"
	"testing"

	"gotest.tools/assert"

	"github.com/go-kit/kit/endpoint"
)

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
			return nil, ErrInvalidCredentials()
		},
		want: "",
		err:  ErrInvalidCredentials(),
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

var endpointCreateTest = []struct {
	name           string
	user           string
	password       string
	createEndpoint endpoint.Endpoint
	err            error
}{
	{
		name:     "should return the error",
		user:     "Jhon",
		password: "pass",
		createEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return CreateResponse{fmt.Errorf("error!")}, nil
		},
		err: fmt.Errorf("error!"),
	},
	{
		name:     "should return an error if the response has the wrong structure",
		user:     "Charles",
		password: "pass",
		createEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return fmt.Errorf("error!"), nil
		},
		err: ErrInvalidResponseStructure(),
	},
	{
		name:     "should return an error if the endpoint returns an error",
		user:     "Jhon",
		password: "pass",
		createEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return nil, ErrInvalidCredentials()
		},
		err: ErrInvalidCredentials(),
	},
}

func TestEndpointCreate(t *testing.T) {
	t.Log("EndpointCreate")

	for _, testcase := range endpointCreateTest {
		t.Logf(testcase.name)

		endpointMock := Endpoints{
			CreateEndpoint: testcase.createEndpoint,
		}
		err := endpointMock.Create(context.Background(), testcase.user, testcase.password)

		var ok bool
		if testcase.err != nil {
			ok = err.Error() == testcase.err.Error()
		} else {
			ok = err == nil
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
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
		name:  "should return the user in the token",
		token: "jjj.www.ttt",
		validateEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return &ValidateResponse{"Jhon", nil}, nil
		},
		want: "Jhon",
	},
	{
		name:  "should return an error if the response has the wrong structure",
		token: "jjj.www.ttt",
		validateEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return "Jhon", nil
		},
		want: "",
		err:  ErrInvalidResponseStructure(),
	},
	{
		name:  "should return an error if the endpoint returns an error",
		token: "jjj.www.ttt",
		validateEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return nil, ErrUserNotFound()
		},
		want: "",
		err:  ErrUserNotFound(),
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

type mockCorrectClientsService struct{}

func (m mockCorrectClientsService) Authorize(ctx context.Context, user, password string) (string, error) {
	return "jjj.www.ttt", nil
}

func (m mockCorrectClientsService) Validate(ctx context.Context, token string) (string, error) {
	return "Jhon", nil
}

func (m mockCorrectClientsService) Create(ctx context.Context, user, password string) error {
	return nil
}

type mockErrorClientsService struct{}

func (m mockErrorClientsService) Authorize(ctx context.Context, user, password string) (string, error) {
	return "", ErrInvalidCredentials()
}

func (m mockErrorClientsService) Validate(ctx context.Context, token string) (string, error) {
	return "", ErrUserNotFound()
}

func (m mockErrorClientsService) Create(ctx context.Context, user, password string) error {
	return ErrInvalidCredentials()
}

var makeAuthorizeEndpointTest = []struct {
	name    string
	client  ClientsService
	request interface{}
	want    *AuthorizeResponse
	err     error
}{
	{
		name:    "should return the token",
		client:  mockCorrectClientsService{},
		request: &AuthorizeRequest{},
		want:    &AuthorizeResponse{"jjj.www.ttt", nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		client:  mockCorrectClientsService{},
		request: "Jhon",
		want:    &AuthorizeResponse{},
		err:     ErrInvalidRequestStructure(),
	},
	{
		name:    "should return an error if the endpoint returns an error",
		client:  mockErrorClientsService{},
		request: &AuthorizeRequest{},
		want:    &AuthorizeResponse{"", ErrInvalidCredentials()},
	},
}

func TestMakeAuthorizeEndpoint(t *testing.T) {
	t.Log("MakeAuthorizeEndpoint")

	for _, testcase := range makeAuthorizeEndpointTest {
		t.Logf(testcase.name)

		endpoint := MakeAuthorizeEndpoint(testcase.client)
		result, err := endpoint(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var makeValidateEndpointTest = []struct {
	name    string
	client  ClientsService
	request interface{}
	want    *ValidateResponse
	err     error
}{
	{
		name:    "should return the user",
		client:  mockCorrectClientsService{},
		request: &ValidateRequest{},
		want:    &ValidateResponse{"Jhon", nil},
	},
	{
		name:    "should return an error if the request has the wrong structure",
		client:  mockCorrectClientsService{},
		request: "jjj.www.ttt",
		want:    &ValidateResponse{},
		err:     ErrInvalidRequestStructure(),
	},
	{
		name:    "should return an error if the endpoint returns an error",
		client:  mockErrorClientsService{},
		request: &ValidateRequest{},
		want:    &ValidateResponse{"", ErrUserNotFound()},
	},
}

func TestMakeValidateEndpoint(t *testing.T) {
	t.Log("MakeValidateEndpoint")

	for _, testcase := range makeValidateEndpointTest {
		t.Logf(testcase.name)

		endpoint := MakeValidateEndpoint(testcase.client)
		result, err := endpoint(context.Background(), testcase.request)

		assert.DeepEqual(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}
