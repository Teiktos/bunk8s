FROM golang:1.16.5-alpine

ENV GRPC_TRACE=all \
    GRPC_VERBOSITY=DEBUG

WORKDIR /build 

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main

WORKDIR /dist

RUN cp /build/main .

EXPOSE 80

CMD ["/dist/main"]