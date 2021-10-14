FROM golang:1.16.5-alpine

ENV GRPC_TRACE=all \
    GRPC_VERBOSITY=DEBUG

ARG CONFIGFILE

WORKDIR /build 

COPY launcher/src/main/go.mod .
COPY launcher/src/main/go.sum .

RUN go mod download

COPY launcher/src/main/. .

RUN go build -o main

WORKDIR /dist

RUN cp /build/main .

COPY ./configuration/$CONFIGFILE /config.yaml

COPY ./launcher/src/cert /cert

CMD ["/dist/main"]
