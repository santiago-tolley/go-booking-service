package commons

import "time"

const (
	ServerHttpAddress = ":8080"
	RoomsGrpcAddr     = ":8081"
	ClientsGrpcAddr   = ":8082"

	JWTSecret     = "very_secret"
	JWTExpiration = 1 * time.Minute
)
