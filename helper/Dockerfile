FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o helper .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/helper .
ENTRYPOINT ["./helper"]
