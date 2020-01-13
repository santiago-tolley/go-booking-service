package rooms

type ErrInvalidResponseStructure struct{}

func (e ErrInvalidResponseStructure) Error() string {
	return "Invalid response structure"
}

type ErrInvalidRequestStructure struct{}

func (e ErrInvalidRequestStructure) Error() string {
	return "Invalid request structure"
}

type ErrNoRoomAvailable struct{}

func (e ErrNoRoomAvailable) Error() string {
	return "No room available"
}
