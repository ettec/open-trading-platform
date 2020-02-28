FROM golang:alpine

RUN apk update && apk upgrade
RUN apk add --no-cache git

RUN go get github.com/fullstorydev/grpcurl
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl


