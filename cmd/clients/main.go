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
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {
	grpcAddr := commons.ClientsGrpcAddr

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	errLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))

	var (
		service = clients.NewClientsServer(clients.WithEncoder(token.JWTEncoder{}),
			clients.WithMongoDB(commons.MongoURL, commons.MongoClientDB))
		endpoints  = clients.MakeEndpoints(service)
		grpcServer = clients.NewGRPCServer(endpoints)
	)

	var g group.Group

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		errLogger.Log("message", "could not set up gRPC listner", "error", err)
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

	logger.Log("gRPC", "listening", "addr", grpcAddr)
	g.Run()
}
