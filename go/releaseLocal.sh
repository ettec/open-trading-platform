#!/bin/bash

echo $PWD
cd $PWD

mkdir -p resources

go build .

if [ $? -eq 0 ]
then
  echo "Successfully ran go build"
else
  echo "go build failed" 
  exit 1
fi



COMPNAME=$(basename "$PWD")

echo Built $COMPNAME


cat > DockerfileLocal << EOF
FROM golang:1.13
ADD $COMPNAME /
COPY resources /resources
CMD /$COMPNAME
EOF

TAG=localhost:32000/$COMPNAME-latest
docker build -f DockerfileLocal -t $TAG .
docker push $TAG


rm $COMPNAME
rm DockerfileLocal


