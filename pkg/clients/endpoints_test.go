package clients

import (
	"context"
	"testing"

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
			return AuthorizeResponse{"jjj.www.ttt", nil}, nil
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
		err:  ErrInvalidResponseStructure{},
	},
	{
		name:     "should return an error if the endpoint returns an error",
		user:     "Jhon",
		password: "pass",
		authorizeEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return nil, ErrInvalidCredentials{}
		},
		want: "",
		err:  ErrInvalidCredentials{},
	},
}

func TestEndpointAuthorize(t *testing.T) {
	t.Log("Authorize")

	for _, testcase := range endpointAuthorizeTest {
		t.Logf(testcase.name)

		endpointMock := Endpoints{
			AuthorizeEndpoint: testcase.authorizeEndpoint,
		}
		result, err := endpointMock.Authorize(context.Background(), testcase.user, testcase.password)

		if !((result != "" && testcase.want != "") || (result == testcase.want)) {
			t.Errorf("=> Got %v wanted %v", result, testcase.want)
		}

		var ok bool
		switch testcase.err.(type) {
		case nil:
			if err == nil {
				ok = true
			}
		case ErrInvalidResponseStructure:
			_, ok = err.(ErrInvalidResponseStructure)
		case ErrInvalidCredentials:
			_, ok = err.(ErrInvalidCredentials)
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
			return ValidateResponse{"Jhon", nil}, nil
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
		err:  ErrInvalidResponseStructure{},
	},
	{
		name:  "should return an error if the endpoint returns an error",
		token: "jjj.www.ttt",
		validateEndpoint: func(_ context.Context, _ interface{}) (interface{}, error) {
			return nil, ErrUserNotFound{}
		},
		want: "",
		err:  ErrUserNotFound{},
	},
}

func TestEndpointValidate(t *testing.T) {
	t.Log("Validate")

	for _, testcase := range endpointValidateTest {
		t.Logf(testcase.name)
		endpointMock := Endpoints{
			ValidateEndpoint: testcase.validateEndpoint,
		}

		result, err := endpointMock.Validate(context.Background(), testcase.token)

		if result != testcase.want {
			t.Errorf("=> Got %v wanted %v", result, testcase.want)
		}

		var ok bool
		switch testcase.err.(type) {
		case nil:
			if err == nil {
				ok = true
			}
		case ErrInvalidToken:
			_, ok = err.(ErrInvalidToken)
		case ErrInvalidResponseStructure:
			_, ok = err.(ErrInvalidResponseStructure)
		case ErrUserNotFound:
			_, ok = err.(ErrUserNotFound)
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
	}
}
