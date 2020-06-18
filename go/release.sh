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


cat > Dockerfile << EOF
FROM ubuntu:19.10
ADD $COMPNAME /
COPY resources /resources
CMD /$COMPNAME
EOF

BASEDIR=$(dirname "$0")

$BASEDIR/../build/pushToDocker.sh

rm $COMPNAME
rm Dockerfile
