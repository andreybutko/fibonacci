package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
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

	ch := make(chan int64, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			num, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			}

			select {
			case ch <- num.Message:
				fmt.Printf("Number recieved %v", num.Message)
				break
			default:
			}
		}
	 close(ch)
	}()

	go func() {
		defer wg.Done()

		for {
			 time.Sleep(100 * time.Millisecond)

			select {
			case num, ok := <-ch:
				if !ok {
					log.Fatalln("Not OK")
					return
				}
				fmt.Printf("%v\n", num)
				break
			}
		}
	}()

	wg.Wait()
}
