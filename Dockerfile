FROM golang:alpine as builder

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git

RUN mkdir /csv-reader
RUN mkdir -p /csv-reader/proto
WORKDIR /csv-reader

COPY ./proto/service.pb.go /csv-reader/proto
COPY ./main.go /csv-reader

COPY go.mod .
COPY go.sum .

#RUN go mod download

RUN go build -o csv-reader .

CMD ./csv-reader