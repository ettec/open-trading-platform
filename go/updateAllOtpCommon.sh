DIRECTORY=$(cd `dirname $0` && pwd)
cd $DIRECTORY
echo working dir `pwd`  

find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line
   echo updating otp common in $line
   go get github.com/ettec/otp-common
   cd $DIRECTORY 
done

echo All updated 


