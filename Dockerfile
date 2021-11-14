FROM golang:alpine

WORKDIR /golang-docker

ADD . .

RUN go mod download

ENTRYPOINT go build  && ./main