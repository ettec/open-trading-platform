FROM golang:1.13

ADD . /app

WORKDIR /app

RUN go build -o service
RUN go test ./...

CMD /app/service
