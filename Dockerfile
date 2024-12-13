FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.18

RUN apk add --no-cache libc6-compat

RUN adduser -D appuser
USER appuser

# Copy the binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Set the command to run the application
CMD ["main", "-migrate"]
