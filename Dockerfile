# Use the official Golang 1.23 image as the base image
FROM golang:1.23.0-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Navigate to the main.go location and build the application
RUN cd cmd/app && go build -o /ticketapp

# Expose the port on which your app runs
EXPOSE 8080

# Command to run the Go app
CMD ["/ticketapp"]
