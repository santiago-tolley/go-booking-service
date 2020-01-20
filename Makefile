build_proto: 
	protoc pb/clients.proto --go_out=plugins=grpc:.
	protoc pb/rooms.proto --go_out=plugins=grpc:.

init_mod: go mod init go-booking-service

build: 
	go build cmd/clients/main.go
	go build cmd/rooms/main.go
	go build cmd/server/main.go

init_db:
	mongo init-mongo.js

docker_init_db:
	docker exec booking_db mongo init-mongo.js

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

docker-start:
	docker-compose up

docker_test:
	docker exec booking_client go test -cover ./pkg/clients/
	docker exec booking_room go test -cover ./pkg/rooms/
	docker exec booking_server go test -cover ./pkg/server/
	docker exec booking_server go test -cover ./pkg/token/
	
docker-test-v:
	docker exec booking_client go test -v -cover ./pkg/clients/
	docker exec booking_room go test -v -cover ./pkg/rooms/
	docker exec booking_server go test -v -cover ./pkg/server/
	docker exec booking_server go test -v -cover ./pkg/token/

docker-build:
	docker build --build-arg service=clients -t clients_service:v1 .
	docker build --build-arg service=rooms -t rooms_service:v1 .
	docker build --build-arg service=server -t server_service:v1 .

docker-db-build:
	docker build -f mongo.Dockerfile -t mongo-image:v1 .

docker-server-ip:
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' booking_server