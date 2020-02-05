FROM golang:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go get -u

RUN go build -o main .


CMD ["/app/main"]