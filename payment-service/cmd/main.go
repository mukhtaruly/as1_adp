package main

import (
	"database/sql"
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"google.golang.org/grpc"

	"payment-service/internal/repository/postgres"
	grpcTransport "payment-service/internal/transport/grpc"
	"payment-service/internal/usecase"
	pb "payment-service/pkg/payment"
)

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Method: %s | Duration: %v | Error: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func main() {
	connStr := os.Getenv("PAYMENT_DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:1234@localhost:5433/payments_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewPaymentRepo(db)
	paymentUsecase := usecase.NewPaymentUsecase(repo)

	
	grpcAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":50051"
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	pb.RegisterPaymentServiceServer(grpcServer, grpcTransport.NewPaymentServer(paymentUsecase))

	go func() {
		log.Println("gRPC Payment Service running on", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to serve gRPC:", err)
		}
	}()

	// ---------- HTTP SERVER (можешь оставить или убрать) ----------
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "payment service ok",
		})
	})

	log.Println("HTTP Payment Service running on :8082")
	r.Run(":8082")
}
