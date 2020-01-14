package rooms

import (
	"context"
	jwt "go-booking-service/pkg/token"
	"sync"
	"testing"
	"time"
)

type validatorCorrect struct{}

func (v validatorCorrect) Validate(_ context.Context, token string) (string, error) {
	return "John", nil
}

type validatorIncorrect struct{}

func (v validatorIncorrect) Validate(_ context.Context, token string) (string, error) {
	return "", jwt.ErrInvalidToken()
}

var serviceBookTest = []struct {
	name      string
	token     string
	date      time.Time
	rooms     []Room
	validator Validator
	want      int
	err       error
}{
	{
		name:  "should return booked room id",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
		},
		validator: validatorCorrect{},
		want:      1,
	},
	{
		name:      "should return en error if there are no rooms available",
		token:     "jjj.www.ttt",
		date:      time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms:     []Room{},
		validator: validatorCorrect{},
		err:       ErrNoRoomAvailable(),
	},
	{
		name:  "should return en error if the token is invalid",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
		},
		validator: validatorIncorrect{},
		err:       jwt.ErrInvalidToken(),
	},
}

func TestServiceBook(t *testing.T) {
	t.Log("ServiceBook")

	for _, testcase := range serviceBookTest {
		t.Logf(testcase.name)

		rs := roomsService{testcase.rooms, testcase.validator}
		result, err := rs.Book(context.Background(), testcase.token, testcase.date)

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

var serviceCheckTest = []struct {
	name      string
	token     string
	date      time.Time
	rooms     []Room
	validator Validator
	want      int
	err       error
}{
	{
		name:  "should retrun the number of available rooms (3)",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
		},
		validator: validatorCorrect{},
		want:      3,
	},
	{
		name:  "should retrun the number of available rooms (0)",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{
					time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC): "John",
				},
				&sync.Mutex{},
			},
		},
		validator: validatorCorrect{},
		want:      0,
	},
}

func TestServiceCheck(t *testing.T) {
	t.Log("ServiceCheck")

	for _, testcase := range serviceCheckTest {
		t.Logf(testcase.name)

		rs := roomsService{testcase.rooms, testcase.validator}
		result, err := rs.Check(context.Background(), testcase.date)

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
