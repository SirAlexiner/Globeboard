# syntax=docker/dockerfile:1.2

# Define the base image and setup build environment
FROM golang:1.22 AS builder

# Define Authors
LABEL authors="Torgrim Thorsen"

# Set the working directory
WORKDIR /app

# Copy go modules for dependencies
COPY go.mod go.sum ./

# Copy the source code.
COPY ../ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/globeboard

# Use a minimal alpine image for the final build stage
FROM alpine:3.19

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose the port on which the application will listen
# Note: This does not publish the port, but documents that the container will listen on this port
EXPOSE 8080

# Define the command to run the application
CMD ["./main"]