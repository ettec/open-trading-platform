#!/bin/bash

# Release all components of the open trading platform to the given docker repo with the given tag

if [ $# -ne 2 ]; then
    echo "Usage: $0 <docker repo> <tag>"
    exit 1
fi

# Assign the parameters to environment variables
export REPO="$1"
export TAG="$2"

echo "Docker Respository: $REPO"
echo "Tag: $TAG"



DIRECTORY=$(cd `dirname $0` && pwd)/..
echo working dir `pwd`	


cd $DIRECTORY
find . -type f -name 'Dockerfile' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line
COMPNAME=otp-$(basename "$PWD")

       echo releasing $COMPNAME
       docker build -t $REPO/$COMPNAME:$TAG .
       docker push $REPO/$COMPNAME:$TAG

cd $DIRECTORY
done

