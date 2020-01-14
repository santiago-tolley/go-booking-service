package server

import "time"

type BookRequest struct {
	Token string    `json:"token"`
	Date  time.Time `json:"date"`
}

type BookResponse struct {
	Id  int   `json:"id"`
	Err error `json:"err"`
}

type CheckRequest struct {
	Date time.Time `json:"date"`
}

type CheckResponse struct {
	Available int   `json:"available"`
	Err       error `json:"err"`
}

type AuthorizeRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type AuthorizeResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	User string `json:"user"`
	Err  error  `json:"err"`
}

func (r AuthorizeResponse) Failed() error {
	return r.Err
}

func (r ValidateResponse) Failed() error {
	return r.Err
}

func (r BookResponse) Failed() error {
	return r.Err
}

func (r CheckResponse) Failed() error {
	return r.Err
}
