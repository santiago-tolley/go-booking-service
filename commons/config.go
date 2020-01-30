package commons

import (
	"os"
	"time"

	"github.com/go-kit/kit/log/level"
)

var (
	ServerHttpAddress = getEnv("SERVER_ADDRESS", ":8080")
	RoomsGrpcAddr     = getEnv("ROOMS_ADDRESS", ":8081")
	ClientsGrpcAddr   = getEnv("CLIENTS_ADDRESS", ":8082")

	JWTSecret     = getEnv("JWT_SECRET", "very_secret")
	JWTExpiration = 2 * time.Hour

	RoomsNumber = 5

	MongoClientURL        = getEnv("MONGO_CLIENT_URL", "mongodb://clients-service:clients-service@localhost:27017/clients-service")
	MongoClientDB         = getEnv("MONGO_CLIENT_DB", "clients-service")
	MongoClientCollection = getEnv("MONGO_CLIENT_COLLECTION", "users")
	MongoClientDBTest     = getEnv("MONGO_CLIENT_DB_TEST", "test-clients-service")

	MongoRoomURL        = getEnv("MONGO_ROOM_URL", "mongodb://rooms-service:rooms-service@localhost:27017/rooms-service")
	MongoRoomDB         = getEnv("MONGO_ROOM_DB", "rooms-service")
	MongoRoomCollection = getEnv("MONGO_ROOM_COLLECTION", "rooms")
	MongoRoomDBTest     = getEnv("MONGO_ROOM_DB_TEST", "test-rooms-service")

	LoggingLevel = level.AllowDebug() // Error > Warn > Info > Debug
)

func getEnv(name, def string) string {
	r := os.Getenv(name)
	if r == "" {
		return def
	}
	return r
}

type contextKey string

func (c contextKey) String() string {
	return "server context key " + string(c)
}

var (
	ContextKeyCorrelationID = contextKey("uuid")
)
