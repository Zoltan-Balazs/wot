FROM golang:alpine

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/wot ./src

ENTRYPOINT ["wot"]
