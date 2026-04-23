package usecase

import (
	"context"

	"payment-service/internal/domain"
	"payment-service/internal/repository"
)

type PaymentUsecase struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: repo}
}


func (u *PaymentUsecase) SavePayment(orderID string, amount float64, status string) error {
	return u.repo.Save(orderID, amount, status)
}


func (u *PaymentUsecase) GetStats(ctx context.Context) (*domain.PaymentStats, error) {
	_ = ctx

	total, paid, failed, totalAmount, err := u.repo.GetStats()
	if err != nil {
		return nil, err
	}

	return &domain.PaymentStats{
		TotalCount:      total,
		AuthorizedCount: paid,
		DeclinedCount:   failed,
		TotalAmount:     totalAmount,
	}, nil
}