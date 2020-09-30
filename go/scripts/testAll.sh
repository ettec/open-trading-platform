DIRECTORY=$(cd `dirname $0` && pwd)/..
cd $DIRECTORY
echo working dir `pwd`	

find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line


   if go test ./...; then
    echo tested $line 
   else
    echo $line tests failed 
    exit 1
   fi 

cd $DIRECTORY
   
done
