# Dockerfile
FROM golang:1.21-alpine

# Install git
RUN apk update && apk add --no-cache git

# Install swag
RUN GO111MODULE=on go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
CMD ["go", "run", "main.go"]