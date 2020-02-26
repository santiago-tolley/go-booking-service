package rooms

import (
	"context"
	"go-booking-service/commons"
	jwt "go-booking-service/pkg/token"
	"os"
	"sync"
	"testing"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	gocmp "github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

var testDB *mongo.Database

func init() {
	errLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	errLogger = kitlog.With(errLogger, "origin", "Test", "caller", kitlog.DefaultCaller)

	testClient, err := mongo.NewClient(options.Client().ApplyURI(commons.MongoClientURL))
	if err != nil {
		errLogger.Log("message", "could not set up mongo client", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = testClient.Connect(ctx)
	if err != nil {
		errLogger.Log("message", "could not connect to database", "error", err)
	}

	testDB = testClient.Database(commons.MongoRoomDBTest)
}

type validatorCorrect struct{}

func (v validatorCorrect) Validate(_ context.Context, token string) (string, error) {
	return "John", nil
}

type validatorIncorrect struct{}

func (v validatorIncorrect) Validate(_ context.Context, token string) (string, error) {
	return "", jwt.ErrInvalidToken()
}

var serviceBookTest = []struct {
	name      string
	token     string
	date      time.Time
	rooms     []Room
	validator Validator
	want      int
	err       error
	init      func(*mongo.Database)
	restore   func(*mongo.Database)
}{
	{
		name:  "should return booked room id",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
		},
		validator: validatorCorrect{},
		want:      0,
	},
	{
		name:      "should return en error if there are no rooms available",
		token:     "jjj.www.ttt",
		date:      time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms:     []Room{},
		validator: validatorCorrect{},
		err:       ErrNoRoomAvailable(),
	},
	{
		name:  "should return en error if the token is invalid",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
		},
		validator: validatorIncorrect{},
		err:       jwt.ErrInvalidToken(),
	},
}

func TestServiceBook(t *testing.T) {
	t.Log("ServiceBook")

	for _, testcase := range serviceBookTest {
		t.Logf(testcase.name)

		rs := &RoomsService{testDB, &testcase.rooms, testcase.validator}
		result, err := rs.Book(context.Background(), testcase.token, testcase.date)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

var serviceCheckTest = []struct {
	name      string
	token     string
	date      time.Time
	rooms     []Room
	validator Validator
	want      int
	err       error
}{
	{
		name:  "should retrun the number of available rooms (3)",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
			{
				map[time.Time]string{},
				&sync.Mutex{},
			},
		},
		validator: validatorCorrect{},
		want:      3,
	},
	{
		name:  "should retrun the number of available rooms (0)",
		token: "jjj.www.ttt",
		date:  time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
		rooms: []Room{
			{
				map[time.Time]string{
					time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC): "John",
				},
				&sync.Mutex{},
			},
		},
		validator: validatorCorrect{},
		want:      0,
	},
}

func TestServiceCheck(t *testing.T) {
	t.Log("ServiceCheck")

	for _, testcase := range serviceCheckTest {
		t.Logf(testcase.name)

		rs := &RoomsService{testDB, &testcase.rooms, testcase.validator}
		result, err := rs.Check(context.Background(), testcase.date)

		assert.Equal(t, result, testcase.want)
		assert.DeepEqual(t, err, testcase.err)
	}
}

type testDecoRoom struct {
	RoomType string        `bson:"type"`
	RoomId   int           `bson:"room_id"`
	Bookings []testBooking `bson:"bookings"`
}

type testBooking struct {
	Date   time.Time `bson:"date"`
	Client string    `bson:"client"`
}

var withRoomsTest = []struct {
	name    string
	rooms   *[]Room
	want    []testDecoRoom
	init    func(*mongo.Database)
	restore func(*mongo.Database)
}{
	{
		name: "Should store the room in the database",
		rooms: &[]Room{
			{
				Book: map[time.Time]string{
					time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC): "John",
				},
			},
		},
		want: []testDecoRoom{
			{
				RoomType: "room",
				RoomId:   0,
				Bookings: []testBooking{
					{
						Date:   time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
						Client: "John",
					},
				},
			},
		},
		init: func(db *mongo.Database) {},
		restore: func(db *mongo.Database) {
			rooms := db.Collection(commons.MongoRoomCollection)
			rooms.Drop(context.Background())
		},
	},
}

func TestWithRooms(t *testing.T) {
	t.Log("WithRooms")
	defer func() {
		rooms := testDB.Collection(commons.MongoRoomCollection)
		rooms.Drop(context.Background())
	}()

	for _, testcase := range withRoomsTest {
		t.Logf(testcase.name)
		testcase.init(testDB)

		rs := &RoomsService{
			db:        testDB,
			validator: validatorCorrect{},
		}

		WithRooms(testcase.rooms)(rs)
		collection := rs.db.Collection(commons.MongoRoomCollection)
		cur, err := collection.Find(context.Background(), bson.D{{"type", "room"}}, nil)
		assert.Assert(t, cmp.Nil(err))

		result := make([]testDecoRoom, 0)
		var room testDecoRoom
		for cur.Next(context.Background()) {
			err = cur.Decode(&room)
			assert.Assert(t, cmp.Nil(err))
			result = append(result, room)
		}
		assert.DeepEqual(t, result, testcase.want)
		testcase.restore(testDB)
	}
}

var withRoomsLoadedTest = []struct {
	name    string
	want    *[]Room
	init    func(*mongo.Database)
	restore func(*mongo.Database)
}{
	{
		name: "Should load the room from the database",
		want: &[]Room{
			{
				Book: map[time.Time]string{
					time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC): "John",
				},
				Mux: &sync.Mutex{},
			},
		},
		init: func(db *mongo.Database) {
			rooms := db.Collection(commons.MongoRoomCollection)
			rooms.InsertOne(context.Background(), bson.M{"type": "meta", "total_rooms": 1})
			rooms.InsertOne(context.Background(), bson.M{
				"type":    "room",
				"room_id": 0,
				"bookings": []bson.M{
					bson.M{
						"date":   time.Date(2020, 6, 13, 12, 0, 0, 0, time.UTC),
						"client": "John",
					},
				},
			})
		},
		restore: func(db *mongo.Database) {
			rooms := db.Collection(commons.MongoRoomCollection)
			rooms.Drop(context.Background())
		},
	},
}

func TestWithLoadedRooms(t *testing.T) {
	t.Log("WithLoadedRooms")
	defer func() {
		rooms := testDB.Collection(commons.MongoRoomCollection)
		rooms.Drop(context.Background())
	}()

	for _, testcase := range withRoomsLoadedTest {
		t.Logf(testcase.name)
		testcase.init(testDB)

		rs := &RoomsService{
			db:        testDB,
			rooms:     &[]Room{},
			validator: validatorCorrect{},
		}
		WithLoadedRooms()(rs)

		assert.DeepEqual(t, rs.rooms, testcase.want, gocmp.AllowUnexported(sync.Mutex{}))
		testcase.restore(testDB)
	}
}
