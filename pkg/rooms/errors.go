package rooms

const (
	InvalidRequestStructure  = "Invalid request structure"
	InvalidResponseStructure = "Invalid response structure"
	NoRoomAvailable          = "No room available"
)

type ErrorWithMsg struct {
	msg string
}

func (e ErrorWithMsg) Error() string {
	return e.msg
}

func ErrInvalidRequestStructure() error {
	return ErrorWithMsg{InvalidRequestStructure}
}

func ErrInvalidResponseStructure() error {
	return ErrorWithMsg{InvalidResponseStructure}
}

func ErrNoRoomAvailable() error {
	return ErrorWithMsg{NoRoomAvailable}
}
