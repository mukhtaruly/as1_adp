package grpc

import (
	"context"

	"payment-service/internal/usecase"
	pb "payment-service/pkg/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	usecase *usecase.PaymentUsecase
}

func NewPaymentServer(u *usecase.PaymentUsecase) *PaymentServer {
	return &PaymentServer{usecase: u}
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	if req.GetOrderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	if req.GetAmount() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}

	message := "Pending"
	success := false

	switch {
	case req.GetAmount() < 500:
		message = "Paid"
		success = true
	case req.GetAmount() <= 5000:
		message = "Pending"
	default:
		message = "Failed"
	}

	// ✅ ВОТ ЭТО НОВОЕ — сохраняем в БД
	err := s.usecase.SavePayment(req.GetOrderId(), req.GetAmount(), message)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PaymentResponse{
		Success: success,
		Message: message,
	}, nil
}

func (s *PaymentServer) GetPaymentStats(ctx context.Context, req *pb.GetPaymentStatsRequest) (*pb.PaymentStats, error) {
	_ = req

	if s.usecase == nil {
		return nil, status.Error(codes.FailedPrecondition, "payment stats usecase is not configured")
	}

	stats, err := s.usecase.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PaymentStats{
		TotalCount:      stats.TotalCount,
		AuthorizedCount: stats.AuthorizedCount,
		DeclinedCount:   stats.DeclinedCount,
		TotalAmount:     stats.TotalAmount,
	}, nil
}
