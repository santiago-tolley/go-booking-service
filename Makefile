build_proto: 
	protoc pb/clients.proto --go_out=plugins=grpc:.
	protoc pb/rooms.proto --go_out=plugins=grpc:.

init_mod: go mod init go-booking-service

build: 
	go build cmd/clients/main.go
	go build cmd/rooms/main.go
	go build cmd/server/main.go

init_db:
	# mongo
	# use clients-service
	# db.createUser({user: "clients-service", pwd: "clients-service", roles: [{role: "readWrite", db: "clients-service"}]})
	# use rooms-service
	# db.createUser({user: "rooms-service", pwd: "rooms-service", roles: [{role: "readWrite", db: "rooms-service"}]})
	mongo clients-service < init-mongo.js

init_db_docker:
	docker exec booking_db mongo clients-service < init-mongo.js
	# dbinit

test:
	go test -cover ./pkg/clients/
	go test -cover ./pkg/rooms/
	go test -cover ./pkg/server/
	go test -cover ./pkg/token/
	
test-v:
	go test -v -cover ./pkg/clients/
	go test -v -cover ./pkg/rooms/
	go test -v -cover ./pkg/server/
	go test -v -cover ./pkg/token/

run_clients:
	go run ./cmd/clients/main.go

run_rooms:
	go run ./cmd/rooms/main.go

run_server:
	go run ./cmd/server/main.go