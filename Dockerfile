# Start with the Go base image
FROM golang:1.18.3-alpine

# Set the working directory
WORKDIR /backend

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o /backend/src/main /backend/src/

# Set the entry point for the container
CMD ["/app/main"]