package commons

import (
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"
)

var (
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))

	ServerHttpAddress = getEnv("SERVER_ADDRESS", "0.0.0.0:8080")
	RoomsGrpcAddr     = getEnv("ROOMS_ADDRESS", "0.0.0.0:8081")
	ClientsGrpcAddr   = getEnv("CLIENTS_ADDRESS", "0.0.0.0:8082")

	JWTSecret     = getEnv("JWT_SECRET", "very_secret")
	JWTExpiration = 10 * time.Minute

	RoomsNumber = 5

	MongoClientURL        = getEnv("MONGO_CLIENT_URL", "mongodb://clients-service:clients-service@localhost:27017/clients-service")
	MongoClientDB         = getEnv("MONGO_CLIENT_DB", "clients-service")
	MongoClientCollection = getEnv("MONGO_CLIENT_COLLECTION", "users")
	MongoClientDBTest     = getEnv("MONGO_CLIENT_DB_TEST", "test-clients-service")

	MongoRoomURL        = getEnv("MONGO_ROOM_URL", "mongodb://rooms-service:rooms-service@localhost:27017/rooms-service")
	MongoRoomDB         = getEnv("MONGO_ROOM_DB", "rooms-service")
	MongoRoomCollection = getEnv("MONGO_ROOM_COLLECTION", "rooms")
	MongoRoomDBTest     = getEnv("MONGO_ROOM_DB_TEST", "test-rooms-service")
)

func getEnv(name, def string) string {
	r := os.Getenv(name)
	if r == "" {
		return def
	}
	return r
}
