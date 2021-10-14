FROM golang:1.16.5-alpine

RUN apk add --no-cache build-base

RUN apk add --no-cache git

WORKDIR /gotest

COPY src/main/go.mod . 
COPY src/main/go.sum . 

RUN go mod download

COPY src/main/. .

RUN go test -c

RUN mkdir /gotest/testresults
