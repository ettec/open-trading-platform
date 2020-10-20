DIRECTORY=$(cd `dirname $0` && pwd)/..
echo working dir `pwd`	


cd $DIRECTORY
find . -type f -name 'Dockerfile' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line
COMPNAME=$(basename "$PWD")

       echo releasing $COMPNAME
       docker build -t localhost:32000/$COMPNAME .
       docker push localhost:32000/$COMPNAME

cd $DIRECTORY
   
done

echo Deployed

