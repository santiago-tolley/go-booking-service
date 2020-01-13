package rooms

import (
	"context"
	"sync"
	"time"
)

type RoomsService interface {
	Book(context.Context, string, time.Time) (int, error)
	Check(context.Context, time.Time) (int, error)
}

type Validator interface {
	Validate(context.Context, string) (string, error)
}

func NewRoomsServer(validator Validator) RoomsService {
	return roomsService{[]room{}, validator}
}

type room struct {
	book map[time.Time]string
	mux  *sync.Mutex
}

type roomsService struct {
	rooms     []room
	validator Validator
}

// Books an availabe room for a date (write/blocking)
// Retruns an error if authentication token is invalid
// or there are no rooms available
func (r roomsService) Book(ctx context.Context, token string, date time.Time) (int, error) {

	// validate token
	user, err := r.validator.Validate(ctx, token)
	if err != nil {
		return 0, err
	}

	var booked bool
	for id, room := range r.rooms {
		if room.book[date] == "" {
			room.mux.Lock()
			if room.book[date] == "" {
				room.book[date] = user
				booked = true
			}
			room.mux.Unlock()
			if booked {
				return id + 1, nil
			}
		}
	}
	return 0, ErrNoRoomAvailable{}
}

// Returns the number of available rooms for a date (read/non-blocking)
func (r roomsService) Check(ctx context.Context, date time.Time) (int, error) {

	var count int
	for _, room := range r.rooms {
		if room.book[date] == "" {
			count++
		}
	}
	return count, nil
}
