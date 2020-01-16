package commons

import "time"

const (
	ServerHttpAddress = ":8080"
	RoomsGrpcAddr     = ":8081"
	ClientsGrpcAddr   = ":8082"

	JWTSecret     = "very_secret"
	JWTExpiration = 10 * time.Minute

	MongoURL              = "mongodb://clients-service:clients-service@localhost:27017/clients-service"
	MongoClientDB         = "clients-service"
	MongoClientCollection = "users"
	MongoClientDBTest     = "test-clients-service"
)
