package token

const (
	InvalidDate      = "Invalid expiration date"
	InvalidAlgorithm = "Invalid signing method"
	InvalidToken     = "Invalid JSON web token"
	ExpiredToken     = "Token is expired"
)

type ErrorWithMsg struct {
	Msg string `json:"message"`
}

func (e ErrorWithMsg) Error() string {
	return e.Msg
}

func ErrInvalidDate() error {
	return ErrorWithMsg{InvalidDate}
}

func ErrInvalidAlgorithm() error {
	return ErrorWithMsg{InvalidAlgorithm}
}

func ErrInvalidToken() error {
	return ErrorWithMsg{InvalidToken}
}

func ErrExpiredToken() error {
	return ErrorWithMsg{ExpiredToken}
}
