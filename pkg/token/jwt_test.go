package token

import (
	"testing"
	"time"
)

var encodeTokenTest = []struct {
	name   string
	user   string
	secret string
	exp    time.Time
	want   string
	err    error
}{
	{
		name:   "should encode token successfuly",
		user:   "Jhon",
		secret: "very_safe",
		exp:    time.Date(2088, 12, 14, 12, 0, 0, 0, time.UTC),
		want:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM3NTM4NjQwMDAsInVzZXIiOiJKaG9uIn0.VZM7zFwJlaBvHNQHAXu-FE30cy8agg2WdvXqygQUGOc",
	},
	{
		name:   "should return an error if the date is in the past",
		user:   "Jhon",
		secret: "very_safe",
		exp:    time.Date(1995, 01, 18, 12, 0, 0, 0, time.UTC),
		want:   "",
		err:    ErrInvalidDate(),
	},
}

func TestEncode(t *testing.T) {
	t.Log("Encode")

	for _, testcase := range encodeTokenTest {
		t.Logf(testcase.name)
		jwtEnc := JWTEncoder{}
		result, err := jwtEnc.Encode(testcase.user, testcase.secret, testcase.exp)

		if result != testcase.want {
			t.Errorf("=> Got %v wanted %v", result, testcase.want)
		}

		if err != testcase.err {
			t.Errorf("=> Got %v wanted %v", err, testcase.err)
		}
	}
}

var decodeTokenTest = []struct {
	name   string
	token  string
	secret string
	want   string
	err    error
}{
	{
		name:   "should return the user in the token",
		token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM3NTM4NjQwMDAsInVzZXIiOiJKb2huIn0.6PFni4sxmV5RYiJ3xRfQUQ7fGrGFB6-hvJtwbBuDLd0",
		secret: "very_safe",
		want:   "John",
	},
	{
		name:   "should return an error if the token is expired",
		token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc5MDQzMDQwMCwidXNlciI6IkpvaG4ifQ.V36DaR7b2FjYer0uaq9KYLDnlJrOgDEo3ayV6lVT1ls",
		secret: "very_safe",
		want:   "",
		err:    ErrExpiredToken(),
	},
}

func TestDecode(t *testing.T) {
	t.Log("Decode")

	for _, testcase := range decodeTokenTest {
		t.Logf(testcase.name)
		jwtDec := JWTEncoder{}
		result, err := jwtDec.Decode(testcase.token, testcase.secret)

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
