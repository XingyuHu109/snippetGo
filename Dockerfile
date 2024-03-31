FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code to the working directory
COPY . .

# Build the Go application
RUN go build -o main ./cmd/web

# Expose the port on which the application will run (adjust if necessary)
EXPOSE 8080

# Set the entry point command to run the application
CMD ["./main"]