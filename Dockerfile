FROM golang:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

ENV TokenTg=1057808441:AAGIh492vskvz82zk7_vbRVyXwDBPUAT9GE
ENV DB_ADDR=postgres://main:cldyvcgaj@192.168.0.107:5431/tossACoin?sslmode\=disable

RUN go get -u

RUN go build -o main .


CMD ["/app/main"]