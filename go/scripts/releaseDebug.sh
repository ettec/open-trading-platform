#!/bin/bash

echo $PWD
cd $PWD



COMPNAME=$(basename "$PWD")



cat > DockerfileDebug << EOF
# Compile stage
FROM golang:1.15.8 AS build-env
# Build Delve
RUN go get github.com/go-delve/delve/cmd/dlv
ADD . /dockerdev
WORKDIR /dockerdev
RUN go build -gcflags="all=-N -l" -o /service
# Final stage
FROM debian:buster
EXPOSE 8000 40000
WORKDIR /
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /service /
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/service"]
EOF

TAG=localhost:32000/$COMPNAME
docker build -f DockerfileDebug -t $TAG .
docker push $TAG


rm DockerfileDebug


