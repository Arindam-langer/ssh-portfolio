# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ssh-portfolio .

# Final runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates openssh-client

# Copy binary and resume
COPY --from=builder /app/ssh-portfolio .
COPY --from=builder /app/resume ./resume
COPY --from=builder /app/arindam\ resume\ latest.pdf .

# Create volume for persistent SSH host keys
VOLUME ["/app/.ssh"]

EXPOSE 2222

ENTRYPOINT ["/app/ssh-portfolio"]
