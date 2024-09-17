# Stage 1: Build the Go application
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy Go modules files and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code to the container
COPY . ./

# Build the Go application
RUN go build -o /my-app

# Stage 2: Run the application in a lightweight container
FROM alpine:latest

WORKDIR /

# Copy the compiled binary from the builder container
COPY --from=builder /my-app /my-app

# Expose the port the app runs on
EXPOSE 8080

# Command to run the app
CMD ["/my-app"]
