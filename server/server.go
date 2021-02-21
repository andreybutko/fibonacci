package main

import (
	"fmt"
	"log"
	"net"

	"github.com/andreybutko/fibonacci/proto"
	"google.golang.org/grpc"
)

// FibonacciServer represents fivonacci server
type FibonacciServer struct {
	proto.UnimplementedFibonacciServer
}

const (
	port = "9002"
)

// GetSequence returns fibonacci sequence stream
func (s *FibonacciServer) GetSequence(_ *proto.FibonacciRequest, stream proto.Fibonacci_GetSequenceServer) error {
	for n := 0; n <= 1000; n++ {
		err := stream.Send(&proto.FibonacciReply{
			Message: int64(n),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	fmt.Println("Server has started.")

	server := grpc.NewServer()
	proto.RegisterFibonacciServer(server, &FibonacciServer{})

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	serverr := server.Serve(lis)
	if serverr != nil {
		log.Fatalf("Server error: %v", serverr)
	}
}
