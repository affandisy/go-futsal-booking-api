# Stage 1: Build the Go application binary
FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/echo-server
RUN go build -o futsal-booking-api main.go

# Stage 2: Create a minimal runtime image and running the application
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY .env .
COPY --from=builder /app/cmd/echo-server/futsal-booking-api /futsal-booking-api
CMD ["/futsal-booking-api"]