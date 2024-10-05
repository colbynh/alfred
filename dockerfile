# Start with a base Go image
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Install Python and pip
RUN apk add --no-cache python3 py3-pip

RUN python3 -m venv /app/kasa
RUN yes | /app/kasa/bin/pip install python-kasa -q -q -q --exists-action i

RUN cp -R /app/kasa/bin/kasa /usr/local/bin/

# Install dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o microservice .

# Expose the port
EXPOSE 8080

CMD ["./microservice"]

# curl -X POST http://localhost:8080/device/192.168.101.170/off