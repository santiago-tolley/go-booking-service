package main

import (
	"fmt"
	"go-booking-service/commons"
	"go-booking-service/pkg/clients"
	"go-booking-service/pkg/rooms"
	"go-booking-service/pkg/server"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {

	httpAddr := commons.ServerHttpAddress
	clientsGrpcAddr := commons.ClientsGrpcAddr
	roomsGrpcAddr := commons.RoomsGrpcAddr

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	errLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))

	clientsGRPCconn, err := grpc.Dial(clientsGrpcAddr, grpc.WithInsecure())
	if err != nil {
		errLogger.Log("transport", "gRPC", "message", "could not connect to clients service", "error", err)
	}

	roomsGRPCconn, err := grpc.Dial(roomsGrpcAddr, grpc.WithInsecure())
	if err != nil {
		errLogger.Log("transport", "gRPC", "message", "could not connect to rooms service", "error", err)
	}

	var (
		service     = server.NewServer(clients.NewGRPCClient(clientsGRPCconn), rooms.NewGRPCClient(roomsGRPCconn))
		endpoints   = server.MakeEndpoints(service)
		httpHandler = server.NewHTTPHandler(endpoints)
	)

	var g group.Group
	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		errLogger.Log("message", "could not set up HTTP listner", "error", err)
	}
	g.Add(func() error {
		return http.Serve(httpListener, httpHandler)
	}, func(error) {
		httpListener.Close()
	})

	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})

	logger.Log("HTTP", "listening", "addr", httpAddr)
	g.Run()
}
