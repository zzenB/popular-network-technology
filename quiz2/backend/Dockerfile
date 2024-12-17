FROM golang:1.23.1 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o pdf-compressor .

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /app/pdf-compressor .
RUN mkdir -p /uploads /compressed
EXPOSE 1337
CMD ["./pdf-compressor"]
