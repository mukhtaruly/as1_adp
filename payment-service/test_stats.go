package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "payment-service/pkg/payment"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := client.GetPaymentStats(ctx, &pb.GetPaymentStatsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TOTAL:", res.TotalCount)
	log.Println("PAID:", res.AuthorizedCount)
	log.Println("FAILED:", res.DeclinedCount)
	log.Println("SUM:", res.TotalAmount)
}