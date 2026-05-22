# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install ca-certificates just in case
RUN apk --no-cache add ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary with CGO disabled for scratch compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bartwarn .

# Runtime stage
FROM scratch

# Copy CA certificates so the app can make HTTPS requests to the BART API and SMTP
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the builder
COPY --from=builder /app/bartwarn /bartwarn

# Run the binary
ENTRYPOINT ["/bartwarn"]
