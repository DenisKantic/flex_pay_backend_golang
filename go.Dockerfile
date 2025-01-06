# Use the official Golang image as the base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire application source code into the container
COPY . .

RUN chmod -R 777 /app

# Build the Go application
RUN go build -o app .
RUN ls -l /app


# Expose the Go application port (8080 in this case)
EXPOSE 8080

# Run the Go application when the container starts
CMD ["./app"]
