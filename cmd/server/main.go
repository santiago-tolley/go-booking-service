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
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, commons.LoggingLevel)
	logger = kitlog.With(logger, "origin", "Server", "caller", kitlog.DefaultCaller)

	clientsGRPCconn, err := grpc.Dial(commons.ClientsGrpcAddr, grpc.WithInsecure())
	if err != nil {
		level.Error(logger).Log("transport", "gRPC", "message", "could not connect to clients service", "error", err)
	}

	roomsGRPCconn, err := grpc.Dial(commons.RoomsGrpcAddr, grpc.WithInsecure())
	if err != nil {
		level.Error(logger).Log("transport", "gRPC", "message", "could not connect to rooms service", "error", err)
	}

	var (
		service     = server.NewServer(clients.NewGRPCClient(clientsGRPCconn), rooms.NewGRPCClient(roomsGRPCconn))
		endpoints   = server.MakeEndpoints(service)
		httpHandler = server.NewHTTPHandler(endpoints)
	)

	var g group.Group
	httpListener, err := net.Listen("tcp", commons.ServerHttpAddress)
	if err != nil {
		level.Error(logger).Log("transport", "HTTP", "message", "could not set up HTTP listner", "error", err)
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

	level.Info(logger).Log("HTTP", "listening", "addr", commons.ServerHttpAddress)
	g.Run()
}
