package main

import (
	"context"
	"fmt"
	"log"
	"io"

	"example.com/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Printf("Created client: %f", c)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber: 3,
		SecondNumber: 5,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculator RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeNumberDecomposition Server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Somehing happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do SquareRoot Unary RPC...")

	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Printf("Error message from server: %v\n", respErr.Message())
			fmt.Printf("Error code from server: %v\n", respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
