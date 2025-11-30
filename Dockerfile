FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o webMessenger .


FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/webMessenger .

EXPOSE 8080
CMD ["./webMessenger"]