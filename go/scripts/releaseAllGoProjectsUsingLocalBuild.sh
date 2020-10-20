DIRECTORY=$(cd `dirname $0` && pwd)/..
echo working dir `pwd`	

RELEASEDFILE=$DIRECTORY/lastreleased.txt


echo -n "" > $RELEASEDFILE 

cd $DIRECTORY
find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line

if ls *.go 1> /dev/null 2>&1;
then
       $DIRECTORY/scripts/releaseGoProjectUsingLocalBuild.sh $RELEASEDFILE 
fi
cd $DIRECTORY
   
done

echo Deployed:
cat $RELEASEDFILE 


