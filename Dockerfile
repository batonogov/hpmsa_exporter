# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o msa_exporter .

# Final stage
FROM alpine:3.22.2

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 10001 exporter && \
    adduser -D -u 10001 -G exporter exporter

WORKDIR /home/exporter

# Copy the binary from builder
COPY --from=builder /build/msa_exporter .

# Change ownership
RUN chown exporter:exporter msa_exporter

# Switch to non-root user
USER exporter

EXPOSE 8000

CMD ./msa_exporter --hostname $HOST --login $LOGIN --password $PASSWORD
