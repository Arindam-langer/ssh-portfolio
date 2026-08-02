# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy all files including vendor directory (100% fast offline build)
COPY . .

# Build lightweight static binary using vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux go build \
    -mod=vendor \
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
COPY --from=builder ["/app/arindam resume latest.pdf", "./arindam resume latest.pdf"]

# Volume for persistent host key generation
VOLUME ["/app/.ssh"]

EXPOSE 2222

ENTRYPOINT ["/app/ssh-portfolio"]
