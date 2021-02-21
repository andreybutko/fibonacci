package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/andreybutko/fibonacci/proto"
)

const (
	address = "localhost:9002"
)

func main() {

	fmt.Println("Client has started.")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFibonacciClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.GetSequence(ctx, &pb.FibonacciRequest{})

	// Skip data when buffer overflowed.
	// TODO: How does it work? Is it awaits til buffer is freed?
	ch := make(chan int64, 2)

	go func() {
		for {
			num, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			}

			ch <- num.Message
		}
		close(ch)
	}()

	// TODO: Why it doesn't work for two goroutines? How to await them?
	for {
		time.Sleep(1 * time.Second)

		num, ok := <-ch

		if !ok {
			break
		}
		fmt.Print(num)
	}
}
