# Start with the official Golang image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files first and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o url-shortener

# Expose the port on which the app will run
EXPOSE 8080

# Command to run the executable
CMD ["./url-shortener"]
