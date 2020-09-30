DIRECTORY=$(cd `dirname $0` && pwd)/..
cd $DIRECTORY
echo working dir `pwd`	

find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line

if ls *.go 1> /dev/null 2>&1;
then

   if go mod tidy; then
    echo mod tidied $line 
   else
    echo failed to mod tidy $line 
    exit 1
   fi 

fi
cd $DIRECTORY
   
done
