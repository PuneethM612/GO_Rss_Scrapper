# Use a minimal Go base image
FROM golang:1.22 AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory
WORKDIR /go/src/github.com/PuneethM06/rssagg

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application binary
RUN go build -o go-rss-scraper .

# Create a lightweight container for running the app
FROM alpine:latest  

WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /go/src/github.com/PuneethM06/rssagg/go-rss-scraper .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./go-rss-scraper"]
