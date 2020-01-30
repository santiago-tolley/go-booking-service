package rooms

const (
	InvalidRequestStructure  = "Invalid request structure"
	InvalidResponseStructure = "Invalid response structure"
	NoRoomAvailable          = "No room available"
	NoCorrelationId          = "No correlation Id in context"
)

type ErrorWithMsg struct {
	Msg string `json:"message"`
}

func (e ErrorWithMsg) Error() string {
	return e.Msg
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

func ErrNoCorrelationId() error {
	return ErrorWithMsg{NoCorrelationId}
}
