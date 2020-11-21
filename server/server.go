package main

import (
	"github.com/andreybutko/fibonacci/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type FibonacciServer struct {
	proto.UnimplementedFibonacciServer
}

const (
	PORT = "9002"
)

func (s *FibonacciServer) GetSequence(_ *proto.FibonacciRequest, stream proto.Fibonacci_GetSequenceServer) error {
	for n := 0; n <= 6; n++ {
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
	server := grpc.NewServer()
	proto.RegisterFibonacciServer(server, &FibonacciServer{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	serverr := server.Serve(lis)
	if serverr != nil {
		log.Fatalf("Server error: %v", serverr)
	}
}
