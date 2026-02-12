# Gunakan versi 1.25 atau yang lebih baru agar sesuai dengan go.mod
FROM golang:1.25-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs
# Install swag binary
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]