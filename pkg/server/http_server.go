package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-booking-service/pkg/clients"
	"go-booking-service/pkg/rooms"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoint Endpoints) http.Handler {

	m := mux.NewRouter()

	m.Methods("POST").Path("/book/{date}").Handler(httptransport.NewServer(
		endpoint.BookEndpoint,
		decodeHTTPBookRequest,
		encodeHTTPGenericResponse,
	))

	m.Methods("GET").Path("/check/{date}").Handler(httptransport.NewServer(
		endpoint.CheckEndpoint,
		decodeHTTPCheckRequest,
		encodeHTTPGenericResponse,
	))

	m.Methods("POST").Path("/authorize/").Handler(httptransport.NewServer(
		endpoint.AuthorizeEndpoint,
		decodeHTTPAuthorizeRequest,
		encodeHTTPGenericResponse,
	))

	m.Methods("POST").Path("/validate/").Handler(httptransport.NewServer(
		endpoint.ValidateEndpoint,
		decodeHTTPValidateRequest,
		encodeHTTPGenericResponse,
	))

	m.Methods("POST").Path("/create/").Handler(httptransport.NewServer(
		endpoint.CreateEndpoint,
		decodeHTTPCreateRequest,
		encodeHTTPGenericResponse,
	))

	return m
}

func decodeHTTPBookRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req = &BookRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, err
	}
	d := mux.Vars(r)["date"]
	date, err := time.Parse("2006-01-02", d)
	req.Date = date
	return req, err
}

func decodeHTTPCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	d := mux.Vars(r)["date"]
	date, err := time.Parse("2006-01-02", d)

	return &CheckRequest{date}, err
}

func decodeHTTPAuthorizeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req = &AuthorizeRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPValidateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req = &ValidateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req = &CreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err.Error() {
	case clients.InvalidRequestStructure:
		return http.StatusBadRequest
	case clients.InvalidResponseStructure:
		return http.StatusBadRequest
	case clients.InvalidCredentials:
		return http.StatusUnauthorized
	case clients.InvalidToken:
		return http.StatusUnauthorized
	case clients.ExpiredToken:
		return http.StatusUnauthorized
	case clients.UserNotFound:
		return http.StatusNotFound
	case rooms.NoRoomAvailable:
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
