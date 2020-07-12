package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"github.com/weilyuwang/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm the client")

	tls := true
	opts := grpc.WithInsecure()

	if tls {
		fmt.Println("With SSL Enabled")
		certFile := "ssl/ca.crt" // Certificate Authority (CA) Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v\n", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close() // defer this line to be executed at the very end

	c := greetpb.NewGreetServiceClient(cc)

	// make requests to grpc server

	// make unary request
	doUnary(c)

	// make server streaming request
	// doServerStreaming(c)

	// make client streaming request
	// doClientStreaming(c)

	// make bi-directional streaming request
	// doBiDirectionalStreaming(c)

	// make unary request with deadline
	// doUnaryWithDeadline(c, 5*time.Second) // should complete
	// doUnaryWithDeadline(c, 1*time.Second) // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a Unary RPC...")

	// contruct client's request

	name := greetpb.Name{
		FirstName: "Alex",
		LastName:  "Spinner",
	}

	req := &greetpb.GreetRequest{
		Greeting: &name,
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Message)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	name := greetpb.Name{
		FirstName: "Jack",
	}

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &name,
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", res.GetMessage())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	// contruct client request
	john := greetpb.Name{
		FirstName: "John",
	}
	stan := greetpb.Name{
		FirstName: "Stan",
	}
	albert := greetpb.Name{
		FirstName: "Albert",
	}
	boris := greetpb.Name{
		FirstName: "Boris",
	}
	jack := greetpb.Name{
		FirstName: "Jack",
	}

	requests := []*greetpb.LongGreatRequest{
		&greetpb.LongGreatRequest{
			Greeting: &john,
		},
		&greetpb.LongGreatRequest{
			Greeting: &stan,
		},
		&greetpb.LongGreatRequest{
			Greeting: &albert,
		},
		&greetpb.LongGreatRequest{
			Greeting: &boris,
		},
		&greetpb.LongGreatRequest{
			Greeting: &jack,
		},
	}

	stream, err := c.LongGreat(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreat: %v", err)
	}

	// We iterate over our requests slice and send each message individually
	// syntax: for index, list[index] := range list
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	// receive response once finish sending all requests
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreat: %v", err)
	}
	fmt.Printf("LongGreet response: %v\n", res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Bi-Directional Streaming Request...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	// contruct client request
	john := greetpb.Name{
		FirstName: "John",
	}
	stan := greetpb.Name{
		FirstName: "Stan",
	}
	albert := greetpb.Name{
		FirstName: "Albert",
	}
	boris := greetpb.Name{
		FirstName: "Boris",
	}
	jack := greetpb.Name{
		FirstName: "Jack",
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &john,
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &stan,
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &albert,
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &boris,
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &jack,
		},
	}

	// make wait channel
	waitc := make(chan struct{})

	// we send a bunch of messages to the server (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		// closesend when we are done stream
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the server (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving message from server stream: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetMessage())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a GreetWithDeadline RPC...")

	// contruct client's request
	name := greetpb.Name{
		FirstName: "John",
		LastName:  "Doe",
	}
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			// grpc error
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline RPC: %v\n", err)
		}
		return
	}

	log.Printf("Response from GreetWithDeadline: %v\n", res.Message)

}
