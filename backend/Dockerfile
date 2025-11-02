#FROM golang:latest AS builder

#WORKDIR /app

# Copy dependency files and download modules
#COPY go.mod go.sum ./
#RUN go mod download

# Copy all source files
#COPY . .

# ✅ Build all Go files, not just main.go
#RUN go build -o ashes-vote-backend .

# ---- Runtime stage ----
#FROM gcr.io/distroless/base-debian12
#WORKDIR /app

#COPY --from=builder /app/ashes-vote-backend /app/
#EXPOSE 8080

#USER nonroot:nonroot
#ENTRYPOINT ["/app/ashes-vote-backend"]



# ---------- Build Stage ----------
FROM golang:latest AS builder
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build a fully static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ashes-vote-backend .

# ---------- Runtime Stage ----------
FROM alpine:3.18
RUN apk add --no-cache ca-certificates bash postgresql-client

WORKDIR /app
COPY --from=builder /app/ashes-vote-backend .

# Optional: Ensure the binary is executable
RUN chmod +x /app/ashes-vote-backend

# Run the app
CMD ["/app/ashes-vote-backend"]


# Use lightweight Linux image with shell for debugging
#FROM alpine:3.18

# Install bash and postgres client for debugging
#RUN apk add --no-cache bash postgresql-client

#WORKDIR /app

# Copy the Go binary
#COPY ./ashes-vote-backend /app/ashes-vote-backend
#RUN chmod +x /app/ashes-vote-backend

# Optional: wait for DB before starting
#COPY wait-for-db.sh /app/wait-for-db.sh
#RUN chmod +x /app/wait-for-db.sh

# Run the wait-for-db script
#CMD ["/app/wait-for-db.sh"]
