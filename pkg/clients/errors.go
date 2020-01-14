package clients

const (
	InvalidCredentials       = "Invalid username or password"
	InvalidToken             = "Invalid JSON Web Token"
	ExpiredToken             = "Expired JSON Web Token"
	UserNotFound             = "User not found"
	InvalidRequestStructure  = "Invalid request structure"
	InvalidResponseStructure = "Invalid response structure"
)

type ErrorWithMsg struct {
	msg string
}

func (e ErrorWithMsg) Error() string {
	return e.msg
}

func ErrInvalidCredentials() error {
	return ErrorWithMsg{InvalidCredentials}
}

func ErrInvalidToken() error {
	return ErrorWithMsg{InvalidToken}
}

func ErrExpiredToken() error {
	return ErrorWithMsg{ExpiredToken}
}

func ErrUserNotFound() error {
	return ErrorWithMsg{UserNotFound}
}

func ErrInvalidRequestStructure() error {
	return ErrorWithMsg{InvalidRequestStructure}
}

func ErrInvalidResponseStructure() error {
	return ErrorWithMsg{InvalidResponseStructure}
}
