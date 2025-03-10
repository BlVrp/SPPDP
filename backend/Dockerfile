# Use the official Golang image as a base
FROM golang:1.23.7-bullseye

# Set the authors label
LABEL authors="volodymyr"

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application files into the container
COPY . .

# Build the Go application
RUN go build -o ./main ./cmd/one-help/main.go

# Expose the port that the Go app will use
EXPOSE 8080

# Set the entry point to run the Go app
CMD ["./main", "run"]
