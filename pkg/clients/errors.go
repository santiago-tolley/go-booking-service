package clients

import "fmt"

type ErrInvalidCredentials struct{}

func (e ErrInvalidCredentials) Error() string {
	return "Invalid username or password"
}

type ErrInvalidToken struct{}

func (e ErrInvalidToken) Error() string {
	return "Invalid JSON Web Token"
}

type ErrUserNotFound struct {
	user string
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("User not found %v", e.user)
}

type ErrInvalidResponseStructure struct{}

func (e ErrInvalidResponseStructure) Error() string {
	return "Invalid response structure"
}

type ErrInvalidRequestStructure struct{}

func (e ErrInvalidRequestStructure) Error() string {
	return "Invalid request structure"
}
