package rooms

import (
	"context"
	"go-booking-service/commons"
	"os"
	"sync"
	"time"

	kitlog "github.com/go-kit/kit/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errLogger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout)) //Stderr

type Service interface {
	Book(context.Context, string, time.Time) (int, error)
	Check(context.Context, time.Time) (int, error)
}

type Validator interface {
	Validate(context.Context, string) (string, error)
}

type ServiceOption func(*RoomsService)

func WithValidator(v Validator) ServiceOption {
	return func(r *RoomsService) {
		r.validator = v
	}
}

func WithRooms(rooms *[]Room) ServiceOption {
	return func(r *RoomsService) {
		if rooms != nil {
			r.rooms = rooms
			if r.db != nil {
				collection := r.db.Collection(commons.MongoRoomCollection)
				_, err := collection.InsertOne(context.Background(), bson.M{"type": "meta", "total_rooms": len(*rooms)})
				if err != nil {
					errLogger.Log("message", "error saving rooms meta", "error", err)
				}

				for i, room := range *rooms {
					bookings := []bson.M{}
					for date, client := range room.Book {
						bookings = append(bookings, bson.M{"date": date, "client": client})
					}
					_, err := collection.InsertOne(context.Background(), bson.M{"type": "room", "room_id": i, "bookings": bookings}) // TODO bulk upload?
					if err != nil {
						errLogger.Log("message", "error saving rooms", "error", err)
					}
				}
			} else {
				errLogger.Log("message", "could not save rooms, no database client set up")
			}
		} else {
			errLogger.Log("message", "room list is empty")
		}
	}
}

func WithLoadedRooms() ServiceOption {
	return func(r *RoomsService) {
		if r.db != nil {
			collection := r.db.Collection(commons.MongoRoomCollection)
			filter := bson.D{{"type", "meta"}}
			var decoMeta struct {
				Type       string `bson:"type"`
				TotalRooms int    `bson:"total_rooms"`
			}
			err := collection.FindOne(context.Background(), filter).Decode(&decoMeta)
			if err != nil {
				errLogger.Log("message", "error loading rooms", "error", err)
				return
			}
			results := make([]Room, decoMeta.TotalRooms)

			cur, err := collection.Find(context.Background(), bson.D{{"type", "room"}}, nil)
			for cur.Next(context.Background()) {

				var decoRoom struct {
					RoomType string `bson:"type"`
					RoomId   int    `bson:"room_id"`
					Bookings []struct {
						Date   time.Time `bson:"date"`
						Client string    `bson:"client"`
					} `bson:"bookings"`
				}
				room := Room{
					Book: map[time.Time]string{},
					Mux:  &sync.Mutex{},
				}
				err = cur.Decode(&decoRoom)
				if err != nil {
					errLogger.Log("message", "error decoding rooms", "error", err)
					return
				}

				for _, b := range decoRoom.Bookings {
					room.Book[b.Date] = b.Client
				}
				results[decoRoom.RoomId] = room
			}

			if err := cur.Err(); err != nil {
				errLogger.Log("message", "error loading rooms", "error", err)
				return
			}
			r.rooms = &results
		} else {
			errLogger.Log("message", "could not load rooms, no database client set up")
		}
	}
}

func WithMongoDB(url, database string) ServiceOption {
	return func(r *RoomsService) {
		db, err := mongo.NewClient(options.Client().ApplyURI(url))
		if err != nil {
			errLogger.Log("message", "could not set up mongo client", "error", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = db.Connect(ctx)
		if err != nil {
			errLogger.Log("message", "could not connect to database", "error", err)
		}
		r.db = db.Database(database)
	}
}

func NewRoomsServer(opts ...ServiceOption) *RoomsService {
	r := &RoomsService{}
	for _, options := range opts {
		options(r)
	}
	return r
}

type Room struct {
	Book map[time.Time]string
	Mux  *sync.Mutex
}

type RoomsService struct {
	db        *mongo.Database
	rooms     *[]Room
	validator Validator
}

// Books an availabe room for a date (write/blocking)
// Retruns an error if authentication token is invalid
// or there are no rooms available
func (r *RoomsService) Book(ctx context.Context, token string, date time.Time) (int, error) {

	// validate token
	user, err := r.validator.Validate(ctx, token)
	if err != nil {
		return 0, err
	}

	var booked bool
	for id, room := range *r.rooms {
		if room.Book[date] == "" {
			room.Mux.Lock()
			if room.Book[date] == "" {
				room.Book[date] = user
				booked = true
			}
			room.Mux.Unlock()
			if booked {
				if r.db != nil {
					users := r.db.Collection(commons.MongoRoomCollection)
					match := bson.M{"room_id": id}
					change := bson.M{"$push": bson.M{"bookings": bson.M{"date": date, "client": user}}} // date format?
					_, err := users.UpdateOne(context.Background(), match, change)
					if err != nil {
						errLogger.Log("message", "error updating room", "error", err)
					}
				} else {
					errLogger.Log("message", "could not save room status, no database client set up")
				}
				return id, nil
			}
		}
	}
	return 0, ErrNoRoomAvailable()
}

// Returns the number of available rooms for a date (read/non-blocking)
func (r *RoomsService) Check(ctx context.Context, date time.Time) (int, error) {

	var count int
	for _, room := range *r.rooms {
		if room.Book[date] == "" {
			count++
		}
	}
	return count, nil
}

func GenerateRooms(total int) []Room {
	r := make([]Room, total)
	for i := 0; i < total; i++ {
		r[i] = Room{
			map[time.Time]string{},
			&sync.Mutex{},
		}
	}
	return r
}
