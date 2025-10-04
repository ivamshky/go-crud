FROM golang:1.22-alpine3.18 AS builder

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix app -o go-crud .


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/go-crud .

EXPOSE 8080

CMD ["./go-crud"]