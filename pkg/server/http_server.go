package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go-booking-service/commons"
	"go-booking-service/pkg/clients"
	"go-booking-service/pkg/rooms"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoint Endpoints) http.Handler {

	m := mux.NewRouter()

	m.Methods("POST").Path("/book/{date}").Handler(httptransport.NewServer(
		endpoint.BookEndpoint,
		decodeHTTPBookRequest,
		encodeHTTPGenericResponse,
		httptransport.ServerBefore(addCorrelationId, validateToken),
	))

	m.Methods("GET").Path("/check/{date}").Handler(httptransport.NewServer(
		endpoint.CheckEndpoint,
		decodeHTTPCheckRequest,
		encodeHTTPGenericResponse,
		httptransport.ServerBefore(addCorrelationId),
	))

	m.Methods("POST").Path("/authorize/").Handler(httptransport.NewServer(
		endpoint.AuthorizeEndpoint,
		decodeHTTPAuthorizeRequest,
		encodeHTTPGenericResponse,
		httptransport.ServerBefore(addCorrelationId),
	))

	m.Methods("POST").Path("/validate/").Handler(httptransport.NewServer(
		endpoint.ValidateEndpoint,
		decodeHTTPValidateRequest,
		encodeHTTPGenericResponse,
		httptransport.ServerBefore(addCorrelationId),
	))

	m.Methods("POST").Path("/create/").Handler(httptransport.NewServer(
		endpoint.CreateEndpoint,
		decodeHTTPCreateRequest,
		encodeHTTPGenericResponse,
		httptransport.ServerBefore(addCorrelationId),
	))

	return m
}

func addCorrelationId(ctx context.Context, r *http.Request) context.Context {
	uuid := uuid.New()

	ctx = context.WithValue(ctx, commons.ContextKeyCorrelationID, uuid)
	return ctx
}

func validateToken(ctx context.Context, r *http.Request) context.Context {
	// take token from request header

	// validate with /validate endpoint

	// decide if request goes through or is blocked... (cancel context? write response?)
	return ctx
}

func decodeHTTPBookRequest(_ context.Context, r *http.Request) (interface{}, error) {
	level.Debug(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding http book request")
	token := getToken(r)
	date, err := getDate(r)
	return &BookRequest{Token: token, Date: date}, err
}

func decodeHTTPCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	level.Debug(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding http check request")
	date, err := getDate(r)
	return &CheckRequest{date}, err
}

func decodeHTTPAuthorizeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	level.Debug(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding http authorize request")
	var req = &AuthorizeRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPValidateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	level.Debug(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding http validate request")
	token := getToken(r)
	return &ValidateRequest{Token: token}, nil
}

func decodeHTTPCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	level.Debug(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "decoding http create request")
	var req = &CreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	level.Debug(logger).Log("correlation ID", ctx.Value(commons.ContextKeyCorrelationID), "message", "encoding http response")
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

// Returns token from header or "" if invalid
func getToken(r *http.Request) string {
	var header = r.Header.Get("Authorization")
	auth_header := strings.Split(header, " ")
	var token string
	if len(auth_header) == 2 {
		token = auth_header[1]
	}
	return token
}

// Returns date from query params or an error if invalid
func getDate(r *http.Request) (time.Time, error) {
	date := mux.Vars(r)["date"]
	return time.Parse("2006-01-02", date)
}
