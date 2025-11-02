# syntax=docker/dockerfile:1
FROM golang:latest AS builder

WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# ✅ Build all Go files, not just main.go
RUN go build -o ashes-vote-backend .

# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /app/ashes-vote-backend /app/
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/ashes-vote-backend"]
