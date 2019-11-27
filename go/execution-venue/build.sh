#!/bin/bash

echo $PWD
cd $PWD

go build .
COMPNAME=$(basename "$PWD")

cat > Dockerfile << EOF
FROM ubuntu:19.10
ADD $COMPNAME /
CMD /$COMPNAME
EOF


LATESTVERSION=$(docker image ls | grep $COMPNAME | head -1 | awk '{print $2}')
LATESTVERSIONIMAGE=$(docker image ls | grep $COMPNAME | head -1 | awk '{print $3}')
NEXTVERSION=$((LATESTVERSION+1))
IMAGEID=$(docker build . | tail -1 | awk '{print $3}')
if [ $IMAGEID = $LATESTVERSIONIMAGE ]  
then
	echo New image id $IMAGEID is identical to the image id of the existing latest version $LATESTVERSION, will exit without tagging and pushing
	exit 1
fi


TAG=192.168.1.200:5000/$COMPNAME:$NEXTVERSION
echo Tagging image id $IMAGEID using tag $TAG:
docker tag $IMAGEID $TAG
docker push $TAG

rm $COMPNAME
rm Dockerfile
