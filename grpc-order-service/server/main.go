package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	orderpb "github.com/neha190495/csk/grpc-order-service/proto"
)

const (
	port = ":50051"
)

type OrderServiceServer struct{}

func main() {
	// Start our listener on the gRPC port: 50054
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Unable to listen on port %s: %v", port, err)
	}

	// Create new gRPC server with (blank) options
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	// Create OrderService type
	srv := &OrderServiceServer{}
	// Register the service with the server
	orderpb.RegisterOrderServiceServer(s, srv)

	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50054")

	// Create a channel to receive OS signals to stop the server using a SHUTDOWN HOOK
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C) and ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	<-c

	// After receiving CTRL+C, properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Done.")

}
