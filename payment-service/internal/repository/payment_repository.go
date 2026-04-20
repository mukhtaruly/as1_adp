package repository

type PaymentRepository interface {
	GetStats() (total int64, paid int64, failed int64, totalAmount int64, err error)
	Save(orderID string, amount float64, status string) error
}