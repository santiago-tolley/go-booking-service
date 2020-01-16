package clients

import (
	"context"
	"go-booking-service/commons"
	jwt "go-booking-service/pkg/token"
	"os"
	"testing"
	"time"

	kitlog "github.com/go-kit/kit/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gotest.tools/assert"
)

var testDB *mongo.Database

func init() {
	errLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))

	testClient, err := mongo.NewClient(options.Client().ApplyURI(commons.MongoURL))
	if err != nil {
		errLogger.Log("message", "could not set up mongo client", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = testClient.Connect(ctx)
	if err != nil {
		errLogger.Log("message", "could not connect to database", "error", err)
	}

	testDB = testClient.Database(commons.MongoClientDBTest)
}

type mockCorrectEncoderDecoder struct{}

func (m mockCorrectEncoderDecoder) Encode(user, secret string, date time.Time) (string, error) {
	return "jjj.www.ttt", nil
}

func (m mockCorrectEncoderDecoder) Decode(token, secret string) (string, error) {
	return "John", nil
}

type mockCorrectNotFoundDecoder struct{}

func (m mockCorrectNotFoundDecoder) Encode(user, secret string, date time.Time) (string, error) {
	return "jjj.www.ttt", nil
}

func (m mockCorrectNotFoundDecoder) Decode(token, secret string) (string, error) {
	return "Charles", nil
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
	encoder  EncoderDecoder
	want     string
	err      error
}{
	{
		name:     "should return the token with the user",
		user:     "John",
		password: "pass",
		encoder:  mockCorrectEncoderDecoder{},
		want:     "jjj.www.ttt",
	},
	{
		name:     "should return an error if the user doesn't exist",
		user:     "Charles",
		password: "pass",
		encoder:  mockCorrectEncoderDecoder{},
		want:     "",
		err:      ErrInvalidCredentials(),
	},
	{
		name:     "should return an error if the password doesn't match",
		user:     "John",
		encoder:  mockCorrectEncoderDecoder{},
		password: "not_pass",
		want:     "",
		err:      ErrInvalidCredentials(),
	},
	{
		name:     "should return an error if the encoder returns an error",
		user:     "John",
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

		c := clientsService{testcase.encoder, testDB}
		result, err := c.Authorize(context.Background(), testcase.user, testcase.password)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var validateTest = []struct {
	name    string
	token   string
	decoder EncoderDecoder
	want    string
	err     error
}{
	{
		name:    "should return the user in the token",
		token:   "jjj.www.ttt",
		decoder: mockCorrectEncoderDecoder{},
		want:    "John",
	},
	{
		name:    "should return an error if the user doesn't exist",
		token:   "jjj.www.ttt",
		decoder: mockCorrectNotFoundDecoder{},
		want:    "",
		err:     ErrUserNotFound(),
	},
	{
		name:    "should return an error if the token is invalid",
		token:   "jjj.www.ttt",
		decoder: mockErrorEncoderDecoder{},
		want:    "",
		err:     jwt.ErrInvalidToken(),
	},
}

func TestValidate(t *testing.T) {
	t.Log("Validate")

	for _, testcase := range validateTest {
		t.Logf(testcase.name)

		c := clientsService{testcase.decoder, testDB}
		result, err := c.Validate(context.Background(), testcase.token)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var createTest = []struct {
	name     string
	user     string
	password string
	err      error
}{
	{
		name:     "should return the error",
		user:     "Charles",
		password: "pass2",
	},
	{
		name:     "should return an error if the user already exists",
		user:     "John",
		password: "pass",
		err:      ErrUserExists(),
	},
}

func TestCreate(t *testing.T) {
	t.Log("Create")

	for _, testcase := range createTest {
		t.Logf(testcase.name)

		c := clientsService{mockCorrectEncoderDecoder{}, testDB}
		err := c.Create(context.Background(), testcase.user, testcase.password)

		assert.DeepEqual(t, err, testcase.err)
	}
}
