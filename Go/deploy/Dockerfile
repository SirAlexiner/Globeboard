# syntax=docker/dockerfile:1
FROM golang:1.22 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY ../go.mod ../go.sum ./

RUN go mod download

# Copy the source code.
COPY ../ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/globeboard

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["./main"]