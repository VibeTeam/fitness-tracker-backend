# -------- Stage 1: builder --------
FROM golang:1.24-alpine AS builder

# Install git (for "go mod download" in private repos, if needed)
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /src

# Copy go module manifests first for better caching
COPY go.mod go.sum ./

# Copy the rest of the workspace (user, shared, workout modules and main package)
COPY . .

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Build the statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/fitness-tracker-backend ./

# -------- Stage 2: runtime --------
FROM gcr.io/distroless/static

COPY --from=builder /app/fitness-tracker-backend /usr/local/bin/fitness-tracker-backend

# Expose HTTP port
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/fitness-tracker-backend"] 