package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"github.com/weilyuwang/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

// Implement server's Greet func
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("Greet function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	greetingMesasge := "Hello " + firstName

	// Contruct server's response
	res := &greetpb.GreetResponse{
		Message: greetingMesasge,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		message := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Message: message,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreat(stream greetpb.GreetService_LongGreatServer) error {

	fmt.Printf("LongGreat function was invoked with a streaming request\n")

	message := ""

	for {
		// read client request from the client stream
		req, err := stream.Recv()

		if err == io.EOF {
			// We have finished reading the client stream
			fmt.Println("Client indicates that all the requests have been sent")
			response := greetpb.LongGreatResponse{
				Message: message,
			}
			return stream.SendAndClose(&response)

		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		message += "Hello " + firstName + "! "
	}

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")

	for {
		// read client request from the client stream
		req, err := stream.Recv()

		if err == io.EOF {
			// The client has done streaming
			fmt.Println("Client is done streaming!")
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		message := "Hello " + firstName + "! "
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Message: message,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to clien: %v", sendErr)
			return sendErr
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {

	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)
	// let server sleep for 3 seconds
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// the client cancelled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	greetingMesasge := "Hello " + firstName

	// Contruct server's response
	res := &greetpb.GreetWithDeadlineResponse{
		Message: greetingMesasge,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello I'm the server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	tls := true
	opts := []grpc.ServerOption{}
	if tls {
		fmt.Println("With SSL Enabled")
		// credentials
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificate: %v\n", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	// Create a grpc server with ssl creds option
	s := grpc.NewServer(opts...)
	// Register the GreetService
	greetpb.RegisterGreetServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
