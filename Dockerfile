# Start from a Go Alpine base image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Install air using go get
RUN go install github.com/cosmtrek/air@latest

# Copy the Go mod and sum files to the working directory
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Expose the port on which the application will listen
EXPOSE 8080

ENTRYPOINT ["air"]
