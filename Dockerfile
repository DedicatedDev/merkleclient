# Use the official Go image to create a build artifact.
FROM golang:1.21.3-bookworm AS builder

# Set the working directory outside $GOPATH to enable Go modules support.
WORKDIR /app

# Retrieve application dependencies.
# This allows for caching of dependencies, improving build speeds.
COPY client/go.mod .
COPY client/go.sum .
RUN go mod download

# Copy the source code as the last step so the build cache can be leveraged
# as much as possible (Go dependencies rarely change, source code often does).
COPY client/ .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./build/client ./main.go

# Use a lightweight Alpine image for the runtime.
FROM alpine:3.14 AS runtime

WORKDIR /app

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/build/client /app/
RUN chmod +x /app/client

COPY client/entrypoint.sh client/1.txt client/2.txt /app/
RUN chmod +x /app/client
RUN chmod +x /app/entrypoint.sh
RUN apk add --no-cache netcat-openbsd
# Set the entrypoint script as the default behavior when the container starts
ENTRYPOINT ["/app/entrypoint.sh"]
#CMD ["upload", "-f", "1.txt", "-f", "2.txt"]
