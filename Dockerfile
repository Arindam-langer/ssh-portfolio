# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy source code including local vendor/
COPY . .

# Build static binary using vendor directory
RUN CGO_ENABLED=0 GOOS=linux go build \
    -mod=mod\
    -ldflags="-s -w" \
    -trimpath \
    -o ssh-portfolio .

# Final minimal runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates openssh-client

# Copy application binary and required assets
COPY --from=builder /app/ssh-portfolio .
COPY --from=builder /app/resume ./resume
COPY --from=builder /app/config.yaml .
# Volume for persistent host key generation
VOLUME ["/app/.ssh"]

EXPOSE 2222

ENTRYPOINT ["/app/ssh-portfolio"]
