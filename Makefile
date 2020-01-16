buildproto: 
	protoc pb/clients.proto --go_out=plugins=grpc:.
	protoc pb/rooms.proto --go_out=plugins=grpc:.

init: go mod init go-booking-service

gobuild: 
	go build cmd/clients/main.go
	go build cmd/rooms/main.go
	go build cmd/server/main.go

dbinit:
	# mongo
	# use clients-service
	# db.createUser({user: "clients-service", pwd: "clients-service", roles: [{role: "readWrite", db: "clients-service"}]})
	# use rooms-service
	# db.createUser({user: "rooms-service", pwd: "rooms-service", roles: [{role: "readWrite", db: "rooms-service"}]})
	mongo clients-service < init-mongo.js

dbinit-docker:
	docker exec booking_db mongo clients-service < init-mongo.js
	# dbinit

test:
	go test -v -cover ./pkg/clients/
	go test -v -cover ./pkg/rooms/
	go test -v -cover ./pkg/server/
	