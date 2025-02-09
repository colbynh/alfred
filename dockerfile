FROM golang:1.23-alpine
WORKDIR /app

RUN apk add --no-cache python3 py3-pip

RUN python3 -m venv /app/kasa
RUN yes | /app/kasa/bin/pip install python-kasa -q -q -q --exists-action i

RUN apk add --no-cache tzdata

RUN cp -R /app/kasa/bin/kasa /usr/local/bin/

COPY go.mod ./
COPY go.sum ./

# Install dependencies
RUN go mod download

COPY . .

WORKDIR /app/cmd/api 

RUN go get .
RUN go build -gcflags=all="-N -l" -o microservice .

EXPOSE 8080

CMD ["./microservice"]

