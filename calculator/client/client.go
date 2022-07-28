package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/i-m-vivek/grpc-assignment/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	Sum(c)

	Prime(c)

	Average(c)

	MaxNum(c)

}

func Sum(c calculatorpb.CalculatorServiceClient) {
	firstnum := 4
	secondnum := 6
	fmt.Printf("\n\n***** Sending request to get sum of %v and %v to server... ***** \n\n", firstnum, secondnum)

	req := calculatorpb.SumRequest{
		FirstNumber:  int32(firstnum),
		SecondNumber: int32(secondnum),
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling Sum: %v", err)
	}
	fmt.Printf("Response: Sum of %v and %v = %v\n", firstnum, secondnum, resp.Result)
}

func Prime(c calculatorpb.CalculatorServiceClient) {
	num := 50
	fmt.Printf("\n\n***** Sending request to get all prime numbers smaller than %v... ***** \n\n", num)

	req := calculatorpb.PrimeRequest{
		Num: int32(num),
	}

	respStream, err := c.Prime(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while Prime server-side streaming grpc : %v", err)
	}
	fmt.Printf("Response: Prime Numbers less than %v => ", num)
	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while receving server stream : %v", err)
		}

		fmt.Printf("%v ", msg.Result)
	}
	fmt.Println()
}

func Average(c calculatorpb.CalculatorServiceClient) {

	fmt.Printf("\n\n***** Sending request to get average of numbers... ***** \n\n")

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing average streaming : %v", err)
	}

	requests := []*calculatorpb.AvgRequest{
		{Num: 1},
		{Num: 2},
		{Num: 3},
		{Num: 4},
		{Num: 5},
		{Num: 6},
	}

	for _, req := range requests {
		stream.Send(req)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}
	fmt.Println("\nResponse: Average =", resp.GetResult())
}

func MaxNum(c calculatorpb.CalculatorServiceClient) {

	requests := []*calculatorpb.MaxNumRequest{
		{Num: 5},
		{Num: 4},
		{Num: 8},
		{Num: 7},
		{Num: 10},
	}
	fmt.Printf("\n\n***** Sending request to get max of numbers...%v\n\n", requests)

	stream, err := c.MaxNum(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming : %v", err)
	}

	waitchan := make(chan struct{})

	go func(requests []*calculatorpb.MaxNumRequest) {
		for _, req := range requests {

			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending request to MaxNum service : %v", err)
			}
		}
		stream.CloseSend()
	}(requests)

	fmt.Printf("Response: ")
	go func() {
		for {

			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				fmt.Println()
				return
			}

			if err != nil {
				log.Fatalf("error receiving response from server : %v", err)
			}

			fmt.Printf("%v ", resp.GetResult())
		}

	}()

	<-waitchan
}
