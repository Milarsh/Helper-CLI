FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o links_store .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/links_store .
EXPOSE 8080
ENTRYPOINT ["./links_store"]
