help:			## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

build:			## Builds the services
	go build cmd/clients/main.go
	go build cmd/rooms/main.go
	go build cmd/server/main.go

build-proto:		## Compiles the protobufs for clients and rooms
	protoc pb/clients.proto --go_out=plugins=grpc:.
	protoc pb/rooms.proto --go_out=plugins=grpc:.

init-db:		## Initializes local mongo db with service users and database indexes
	mongo init-mongo.js

init-mod:		## Initalizes go modules
	go mod init go-booking-service

run-clients:		## Starts the clients service
	go run ./cmd/clients/main.go

run-rooms:		## Starts the rooms service
	go run ./cmd/rooms/main.go

run-server:		## Starts the server service
	go run ./cmd/server/main.go

test:			## Runs service tests
	go test -cover ./pkg/clients/
	go test -cover ./pkg/rooms/
	go test -cover ./pkg/server/
	go test -cover ./pkg/token/
	
test-v:			# Runs service tests with verbose output
	go test -v -cover ./pkg/clients/
	go test -v -cover ./pkg/rooms/
	go test -v -cover ./pkg/server/
	go test -v -cover ./pkg/token/

docker-build:		## Builds the service docker images
	docker build -f clients.Dockerfile -t clients-service:v1 . 
	docker build -f rooms.Dockerfile -t rooms-service:v1 .
	docker build -f server.Dockerfile -t server-service:v1 .

docker-build-db:	## Builds mongo database image
	docker build -f mongo.Dockerfile -t mongo-image:v1 .

docker-init-db:		## Initializes running docker mongo database
	docker exec booking_db mongo init-mongo.js

docker-server-ip:	## Shows the IP address of the running HTTP server instance 
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' booking_server

docker-start:		## Starts mongo database and service containers
	docker-compose up

docker-start-d:		## Starts mongo database and service containers detatched from the console
	docker-compose up -d

docker-test:		## Runs service tests from running docker instances 
	docker exec booking_clients go test -cover ./pkg/clients/
	docker exec booking_rooms go test -cover ./pkg/rooms/
	docker exec booking_server go test -cover ./pkg/server/
	docker exec booking_server go test -cover ./pkg/token/
	
docker-test-v:		## Runs service tests from running docker instances with verbose output
	docker exec booking_clients go test -v -cover ./pkg/clients/
	docker exec booking_rooms go test -v -cover ./pkg/rooms/
	docker exec booking_server go test -v -cover ./pkg/server/
	docker exec booking_server go test -v -cover ./pkg/token/
