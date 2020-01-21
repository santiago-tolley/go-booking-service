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

	JWTSecret     = getEnv("JWTSecret", "very_secret")
	JWTExpiration = 10 * time.Minute

	RoomsNumber = 5

	MongoClientURL        = getEnv("MongoClientURL", "mongodb://clients-service:clients-service@localhost:27017/clients-service")
	MongoClientDB         = getEnv("MongoClientDB", "clients-service")
	MongoClientCollection = getEnv("MongoClientCollection", "users")
	MongoClientDBTest     = getEnv("MongoClientDBTest", "test-clients-service")

	MongoRoomURL        = getEnv("MongoRoomURL", "mongodb://rooms-service:rooms-service@localhost:27017/rooms-service")
	MongoRoomDB         = getEnv("MongoRoomDB", "rooms-service")
	MongoRoomCollection = getEnv("MongoRoomCollection", "rooms")
	MongoRoomDBTest     = getEnv("MongoRoomDBTest", "test-rooms-service")
)

func getEnv(name, def string) string {
	r := os.Getenv(name)
	if r == "" {
		return def
	}
	return r
}
