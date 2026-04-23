package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "order-service/pkg/order"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	stream, err := client.SubscribeToOrderUpdates(ctx, &pb.OrderRequest{
		OrderId: "7fa21429-7a2d-442e-998d-ba58ccea54b4",
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("STATUS:", res.Status)
	}
}
