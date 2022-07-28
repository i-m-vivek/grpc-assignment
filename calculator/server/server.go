package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/i-m-vivek/grpc-assignment/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (resp *calculatorpb.SumResponse, err error) {
	fmt.Println("\n\nGot a request for sum")

	firstnum := req.GetFirstNumber()
	secondnum := req.GetSecondNumber()
	sum := firstnum + secondnum
	resp = &calculatorpb.SumResponse{
		Result: sum,
	}
	return resp, nil
}

// is prime or not
func isPrime(n int) bool {
	// Corner case
	if n <= 1 {
		return false
	}

	// Check from 2 to n-1
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func (*server) Prime(req *calculatorpb.PrimeRequest, resp calculatorpb.CalculatorService_PrimeServer) error {

	fmt.Println("\n\nGot request for prime numbers...")

	num := int(req.GetNum())

	for i := 2; i <= num; i++ {
		if isPrime(i) {
			res := calculatorpb.PrimeResponse{
				Result: int32(i),
			}

			resp.Send(&res)
		}

	}

	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {

	fmt.Println("\n\nGot request for average of stream...")

	sum := 0
	count := 0

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading client stream
			return stream.SendAndClose(&calculatorpb.AvgResponse{
				Result: float32(sum) / float32(count),
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}

		count++
		sum += int(msg.GetNum())

	}
}

func (*server) MaxNum(stream calculatorpb.CalculatorService_MaxNumServer) error {
	fmt.Println("\n\nGot request for max of stream...")
	maxnum := -99999

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while receiving data from MaxNum client : %v", err)
			return err
		}

		num := req.GetNum()

		if num > int32(maxnum) {

			sendErr := stream.Send(&calculatorpb.MaxNumResponse{
				Result: num,
			})

			if sendErr != nil {
				log.Fatalf("error while sending response to MaxNum Client : %v", err)
				return err
			}
			maxnum = int(num)
		}

	}
}

func main() {
	fmt.Println("Starting the server...")
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// //Register reflection service on gRPC server
	// reflection.Register(s)

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
