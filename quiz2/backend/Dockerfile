FROM golang:1.23.1 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pdf-compressor .

FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/pdf-compressor .
RUN mkdir -p /app/uploads /app/downloads
EXPOSE 1337
CMD ["./pdf-compressor"]
