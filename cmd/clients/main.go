package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go-booking-service/commons"
	"go-booking-service/pb"
	"go-booking-service/pkg/clients"
	"go-booking-service/pkg/token"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, commons.LoggingLevel)
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)

	var (
		service = clients.NewClientsServer(clients.WithEncoder(token.JWTEncoder{}),
			clients.WithMongoDB(commons.MongoClientURL, commons.MongoClientDB))
		endpoints  = clients.MakeEndpoints(service)
		grpcServer = clients.NewGRPCServer(endpoints)
	)

	var g group.Group

	grpcListener, err := net.Listen("tcp", commons.ClientsGrpcAddr)
	if err != nil {
		level.Error(logger).Log("transport", "gRPC", "message", "could not set up gRPC listner", "error", err)
	}

	g.Add(func() error {
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		pb.RegisterClientsServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
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

	level.Info(logger).Log("gRPC", "listening", "addr", commons.ClientsGrpcAddr)
	g.Run()
}
