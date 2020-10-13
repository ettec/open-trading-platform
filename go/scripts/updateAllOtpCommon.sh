#Specify the tag as the first argument e.g. @v1.0.1 or leave it blank to update to head of the master branch

DIRECTORY=$(cd `dirname $0` && pwd)/..
cd $DIRECTORY
echo working dir `pwd`  

VERSION=$1

find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line
   echo updating otp common in $line
   go get github.com/ettec/otp-common$VERSION
   cd $DIRECTORY 
done

echo All updated 


