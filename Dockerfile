FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE ${SERVER_PORT}
CMD ["./server"]