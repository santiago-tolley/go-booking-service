package clients

import (
	"context"
	"go-booking-service/pkg/token"
	"testing"
)

var authorizeTest = []struct {
	name     string
	user     string
	password string
	users    map[string]string
	want     string
	err      error
}{
	{
		name:     "should return the token with the user",
		user:     "Jhon",
		password: "pass",
		users: map[string]string{
			"Jhon": "pass",
		},
		want: "Jhon",
	},
	{
		name:     "should return an error if the user doesn't exist",
		user:     "Charles",
		password: "pass",
		users: map[string]string{
			"Jhon": "pass",
		},
		want: "",
		err:  ErrInvalidCredentials{},
	},
	{
		name: "should return an error if the password doesn't match",
		user: "Jhon",
		users: map[string]string{
			"Jhon": "pass",
		},
		password: "not_pass",
		want:     "",
		err:      ErrInvalidCredentials{},
	},
}

func TestAuthorize(t *testing.T) {
	t.Log("Authorize")

	for _, testcase := range authorizeTest {
		t.Logf(testcase.name)

		c := ClientsService{token.JWTEncoder{}, testcase.users}
		result, err := c.Authorize(context.Background(), testcase.user, testcase.password)

		if !((result != "" && testcase.want != "") || (result == testcase.want)) {
			t.Errorf("=> Got %v wanted %v", result, testcase.want)
		}

		var ok bool
		switch testcase.err.(type) {
		case nil:
			if err == nil {
				ok = true
			}
		case ErrInvalidCredentials:
			_, ok = err.(ErrInvalidCredentials)
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
	}
}

var validateTest = []struct {
	name  string
	token string
	users map[string]string
	want  string
	err   error
}{
	{
		name:  "should return the user in the token",
		token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM3NTM4NjQwMDAsInVzZXIiOiJKaG9uIn0.VZM7zFwJlaBvHNQHAXu-FE30cy8agg2WdvXqygQUGOc",
		users: map[string]string{
			"Jhon": "pass",
		},
		want: "Jhon",
	},
	{
		name:  "should return an error if the user doesn't exist",
		token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM3NTM4NjQwMDAsInVzZXIiOiJKaG9uIn0.VZM7zFwJlaBvHNQHAXu-FE30cy8agg2WdvXqygQUGOc",
		users: map[string]string{
			"Charles": "pass",
		},
		want: "",
		err:  ErrUserNotFound{},
	},
}

func TestValidate(t *testing.T) {
	t.Log("Validate")

	for _, testcase := range validateTest {
		t.Logf(testcase.name)

		c := ClientsService{token.JWTEncoder{}, testcase.users}
		result, err := c.Validate(context.Background(), testcase.token)

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
		case ErrUserNotFound:
			_, ok = err.(ErrUserNotFound)
		}
		if !ok {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
	}
}
