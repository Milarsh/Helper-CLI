FROM golang:1.24.3-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o articles .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/articles .
EXPOSE 8081
ENTRYPOINT ["./articles"]
