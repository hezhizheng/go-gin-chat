FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "main.go"]