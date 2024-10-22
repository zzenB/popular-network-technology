FROM golang:1.23.1-alpine AS build

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]
