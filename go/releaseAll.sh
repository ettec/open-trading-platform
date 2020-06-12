DIRECTORY=$(cd `dirname $0` && pwd)
cd $DIRECTORY
echo working dir `pwd`	

RELEASEDFILE=$DIRECTORY/lastreleased.txt


echo -n "" > $RELEASEDFILE 

find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line

if ls *.go 1> /dev/null 2>&1;
then
       $DIRECTORY/release.sh $RELEASEDFILE 
fi
cd $DIRECTORY
   
done

echo Deployed:
cat $RELEASEDFILE 


