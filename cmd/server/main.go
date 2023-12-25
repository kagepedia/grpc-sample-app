package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	hellopb "grpc-sample-app/pb/greet"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type GreetingServiceServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (server *GreetingServiceServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello %s, you are %d years old", req.GetName(), req.GetAge()),
	}, nil
}

func (server *GreetingServiceServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	retryCount := 5
	for i := 0; i < retryCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello %s, you are %d years old", i, req.GetName(), req.GetAge()),
		}); err != nil {
			log.Printf("failed to send message: %v", err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func NewGreetingServiceServer() *GreetingServiceServer {
	return &GreetingServiceServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	// Register the service with the server
	hellopb.RegisterGreetingServiceServer(grpcServer, NewGreetingServiceServer())

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server is running on port %v", port)
		grpcServer.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}
