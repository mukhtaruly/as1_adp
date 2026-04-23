package postgres

import "database/sql"

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Save(orderID string, amount float64, status string) error {
	_, err := r.db.Exec(`
		INSERT INTO payments (order_id, amount, status)
		VALUES ($1, $2, $3)
	`, orderID, amount, status)
	return err
}

func (r *PaymentRepo) GetStats() (total int64, paid int64, failed int64, totalAmount int64, err error) {
	err = r.db.QueryRow(`
		SELECT
			COUNT(*) AS total_count,
			COUNT(*) FILTER (WHERE status = 'Paid') AS paid_count,
			COUNT(*) FILTER (WHERE status != 'Paid') AS failed_count,
			COALESCE(SUM(amount)::bigint, 0) AS total_amount
		FROM payments
	`).Scan(&total, &paid, &failed, &totalAmount)

	return
}