FROM golang:1.13

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build
RUN go test ./...


CMD /app/fix-sim-execution-venue
