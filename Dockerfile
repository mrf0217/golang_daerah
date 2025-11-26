

# ---- Build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install build deps
RUN apk add --no-cache git

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/app

# ---- Run stage ----
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# Set runtime env (can be overridden by compose)
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_SSLMODE=disable

# Copy binary
COPY --from=builder /app/server /app/server

# Expose port
EXPOSE 8080

# Run
CMD ["/app/server"]
