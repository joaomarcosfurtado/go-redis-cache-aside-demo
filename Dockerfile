FROM golang:1.23-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

ENV PORT=8080

CMD ["./main"]