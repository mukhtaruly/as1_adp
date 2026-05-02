FROM golang:1.25

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd

CMD ["./main"]