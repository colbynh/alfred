FROM golang:1.23.5-alpine3.21
WORKDIR /app

# Install system dependencies
RUN apk update && apk add --no-cache \
    python3 \
    py3-pip \
    tzdata \
    curl \
    nmap \
    git

# Setup Python Kasa
RUN python3 -m venv /tmp/kasa && \
    yes | /tmp/kasa/bin/pip install python-kasa -q -q -q --exists-action i && \
    cp -R /tmp/kasa/bin/kasa /usr/local/bin/

# Install Air for live reloading
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Install godoc
RUN go install golang.org/x/tools/cmd/godoc@latest

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose ports for both the API and godoc
EXPOSE 8080 6060

WORKDIR /app

# Start both the API and godoc server
CMD ["sh", "-c", "./bin/air -c .air.toml & godoc -http=:6060"]