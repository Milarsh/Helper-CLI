FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o helper .

FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates dumb-init
WORKDIR /app
COPY --from=builder /app/helper .
COPY cron/root /etc/crontabs/root
RUN touch /var/log/cron.log
ENTRYPOINT ["dumb-init", "--", "/usr/sbin/crond", "-f", "-L", "/dev/stdout"]