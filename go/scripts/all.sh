# Script to make it easy to run go commands across all go projects.  E.g.  "./all.sh go test ./..." will run the tests of all go projects.
# The alternative approach here would be to combine all services into one go project, however on balance it was considered that having each
# distinct service in its own project aids comprehension.

DIRECTORY=$(cd `dirname $0` && pwd)/..
cd $DIRECTORY
echo working dir `pwd`	

cmd="$*"

find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line

if ls *.go 1> /dev/null 2>&1;
then

   if eval "$cmd"; then
    echo executed "$cmd" in "$line"
   else
    echo failed to execute "$cmd" 
    exit 1
   fi 

fi
cd $DIRECTORY
   
done
