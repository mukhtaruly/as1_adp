package usecase

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"

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
	// 1. Сохраняем в БД
	err := u.repo.Save(orderID, amount, status)
	if err != nil {
		return err
	}

	// 2. Публикуем событие (RabbitMQ)
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Println("RabbitMQ connection error:", err)
		return nil // не ломаем основной процесс
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Channel error:", err)
		return nil
	}
	defer ch.Close()

	// объявляем очередь (durable)
	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Queue declare error:", err)
		return nil
	}

	// создаём JSON
	bodyData := map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
		"status":   status,
	}

	body, _ := json.Marshal(bodyData)

	// отправка сообщения
	err = ch.Publish(
		"",
		"payment.completed",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Publish error:", err)
	}

	return nil
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