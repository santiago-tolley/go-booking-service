package rooms

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
