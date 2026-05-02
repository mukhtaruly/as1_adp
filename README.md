Assignment 3: Event-Driven Architecture (EDA) with RabbitMQ

Overview

This project extends the previous gRPC-based microservices system by introducing Event-Driven Architecture (EDA) using RabbitMQ.

The system consists of three microservices:

* Order Service (HTTP + gRPC client)
* Payment Service (gRPC server + RabbitMQ producer)
* Notification Service (RabbitMQ consumer)

⸻

Architecture Diagram

(Добавь сюда свою картинку диаграммы, которую мы сделали)

⸻

Tech Stack

* Golang
* gRPC
* RabbitMQ
* PostgreSQL
* Docker & Docker Compose

⸻

How to Run

docker-compose up –build

⸻

API

Create Order:

curl -X POST http://localhost:8081/orders 
-H “Content-Type: application/json” 
-d ‘{“user_id”:1,“amount”:100}’

⸻

Event Flow

1. Client sends HTTP request to Order Service
2. Order Service saves order to PostgreSQL
3. Order Service calls Payment Service via gRPC
4. Payment Service processes payment
5. Payment Service saves payment to DB (transaction commit)
6. Payment Service publishes event to RabbitMQ (payment.completed)
7. Notification Service consumes the event
8. Notification Service logs notification (simulates email)

⸻

Message Format

{
“order_id”: “uuid”,
“amount”: 100,
“status”: “Paid”
}

⸻

Reliability & Delivery Guarantees

Manual ACK:

* Auto-ack is disabled
* Messages are acknowledged only after successful processing

Durable Queue:

* Queue payment.completed is durable
* Messages survive broker restart

Retry Mechanism:

* Services wait for DB and RabbitMQ before starting

⸻

Idempotency Strategy

To avoid duplicate processing:

* Each message has unique order_id
* Processed messages can be tracked
* Duplicate events are ignored

⸻

Docker Services

* postgres
* rabbitmq
* order-service
* payment-service
* notification-service

⸻

Key Design Principles

* At-least-once delivery
* Loose coupling
* Separation of concerns
* Fault tolerance

⸻

RabbitMQ UI

http://localhost:15672
login: guest / guest

⸻

Health Check

curl http://localhost:8081/health
curl http://localhost:8082/health

⸻

Project Structure

.
├── order-service
├── payment-service
├── notification-service
├── db
├── proto
├── docker-compose.yml
└── README.md

⸻

Result

* Microservices communicate via gRPC and RabbitMQ
* Event-driven flow is implemented
* System runs полностью через Docker

⸻

Conclusion

This project demonstrates the transition from synchronous gRPC communication to asynchronous event-driven architecture using message queues, improving scalability, reliability, and decoupling between services.
