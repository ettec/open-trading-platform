FROM golang:1.13

RUN mkdir /common

ADD ../../common/. common

RUN mkdir /execution-venues/fix-sim-execution-venue

ADD . /execution-venues/fix-sim-execution-venue

WORKDIR /execution-venues/fix-sim-execution-venue

RUN go build
RUN go test ./...


CMD /app/fix-sim-execution-venue
