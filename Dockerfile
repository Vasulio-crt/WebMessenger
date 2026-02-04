FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o webMessenger .


FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/webMessenger .

COPY resource resource
COPY pages pages

EXPOSE 8080
CMD ["./webMessenger"]