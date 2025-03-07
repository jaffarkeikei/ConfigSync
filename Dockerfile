# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /workspace

# Copy Go module files
COPY go.mod go.mod
COPY go.sum go.sum

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o configsync-operator cmd/operator/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /

# Install git, kubectl, and ca-certificates
RUN apk add --no-cache git curl ca-certificates \
    && curl -LO "https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/

# Copy the binary from builder
COPY --from=builder /workspace/configsync-operator .

# Run as non-root
USER 65532:65532

ENTRYPOINT ["/configsync-operator"] 