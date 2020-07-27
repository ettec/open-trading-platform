COMPNAME=$(basename "$PWD")

echo Built $COMPNAME


TAG=localhost:5000/$COMPNAME
docker build -f Dockerfile -t $TAG .
docker push $TAG

