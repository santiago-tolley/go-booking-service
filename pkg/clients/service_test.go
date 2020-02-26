package clients

import (
	"context"
	jwt "go-booking-service/pkg/token"
	"testing"
	"time"

	"gotest.tools/assert"
)

type mockCorrectEncoderDecoder struct{}

func (m mockCorrectEncoderDecoder) Encode(user, secret string, date time.Time) (string, error) {
	return "jjj.www.ttt", nil
}

func (m mockCorrectEncoderDecoder) Decode(token, secret string) (string, error) {
	return "John", nil
}

type mockErrorEncoderDecoder struct{}

func (m mockErrorEncoderDecoder) Encode(user, secret string, date time.Time) (string, error) {
	return "", jwt.ErrInvalidToken()
}

func (m mockErrorEncoderDecoder) Decode(token, secret string) (string, error) {
	return "", jwt.ErrInvalidToken()
}

var authorizeTest = []struct {
	name     string
	user     string
	password string
	users    map[string]string
	encoder  EncoderDecoder
	want     string
	err      error
}{
	{
		name:     "should return the token with the user",
		user:     "John",
		password: "pass",
		users: map[string]string{
			"John": "pass",
		},
		encoder: mockCorrectEncoderDecoder{},
		want:    "jjj.www.ttt",
	},
	{
		name:     "should return an error if the user doesn't exist",
		user:     "Charles",
		password: "pass",
		users: map[string]string{
			"John": "pass",
		},
		encoder: mockCorrectEncoderDecoder{},
		want:    "",
		err:     ErrInvalidCredentials(),
	},
	{
		name: "should return an error if the password doesn't match",
		user: "John",
		users: map[string]string{
			"John": "pass",
		},
		encoder:  mockCorrectEncoderDecoder{},
		password: "not_pass",
		want:     "",
		err:      ErrInvalidCredentials(),
	},
	{
		name: "should return an error if the encoder returns an error",
		user: "John",
		users: map[string]string{
			"John": "pass",
		},
		encoder:  mockErrorEncoderDecoder{},
		password: "pass",
		want:     "",
		err:      jwt.ErrInvalidToken(),
	},
}

func TestAuthorize(t *testing.T) {
	t.Log("Authorize")

	for _, testcase := range authorizeTest {
		t.Logf(testcase.name)

		c := clientsService{testcase.encoder, testcase.users}
		result, err := c.Authorize(context.Background(), testcase.user, testcase.password)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var validateTest = []struct {
	name    string
	token   string
	users   map[string]string
	decoder EncoderDecoder
	want    string
	err     error
}{
	{
		name:  "should return the user in the token",
		token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM3NTM4NjQwMDAsInVzZXIiOiJKaG9uIn0.VZM7zFwJlaBvHNQHAXu-FE30cy8agg2WdvXqygQUGOc",
		users: map[string]string{
			"John": "pass",
		},
		decoder: mockCorrectEncoderDecoder{},
		want:    "John",
	},
	{
		name:  "should return an error if the user doesn't exist",
		token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM3NTM4NjQwMDAsInVzZXIiOiJKaG9uIn0.VZM7zFwJlaBvHNQHAXu-FE30cy8agg2WdvXqygQUGOc",
		users: map[string]string{
			"Charles": "pass",
		},
		decoder: mockCorrectEncoderDecoder{},
		want:    "",
		err:     ErrUserNotFound(),
	},
	{
		name:  "should return an error if the token is invalid",
		token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM3NTM4NjQwMDAsInVzZXIiOiJKaG9uIn0.VZM7zFwJlaBvHNQHAXu-FE30cy8agg2WdvXqygQUGOc",
		users: map[string]string{
			"Charles": "pass",
		},
		decoder: mockErrorEncoderDecoder{},
		want:    "",
		err:     jwt.ErrInvalidToken(),
	},
}

func TestValidate(t *testing.T) {
	t.Log("Validate")

	for _, testcase := range validateTest {
		t.Logf(testcase.name)

		c := clientsService{testcase.decoder, testcase.users}
		result, err := c.Validate(context.Background(), testcase.token)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}
