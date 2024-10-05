# Start with a base Go image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Install Python and pip
RUN apk add --no-cache python3 py3-pip

# Install python-kasa using pip
RUN pip install python-kasa

# Install dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o microservice .

# Expose the port
EXPOSE 8080

# Command to run the executable
CMD ["./microservice"]
