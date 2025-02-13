FROM golang:1.23.5-alpine3.21
WORKDIR /app

RUN apk update
RUN apk add --no-cache python3 py3-pip

RUN python3 -m venv /tmp/kasa
RUN yes | /tmp/kasa/bin/pip install python-kasa -q -q -q --exists-action i

RUN apk add --no-cache tzdata 						      # update the local registry
RUN apk add curl

RUN cp -R /tmp/kasa/bin/kasa /usr/local/bin/

RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/api 

RUN go get .

EXPOSE 8080

WORKDIR /app

CMD ["./bin/air", "-c", ".air.toml"]