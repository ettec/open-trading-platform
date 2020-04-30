DIRECTORY=$(cd `dirname $0` && pwd)
cd $DIRECTORY
echo working dir `pwd`	

DEPLOYEDFILE=$DIRECTORY/lastdeployed.txt


echo -n "" > $DEPLOYEDFILE 

find . -type f -name '*go.mod*' | sed -r 's|/[^/]+$||' |sort |uniq | while read line; do
cd $line

if ls *.go 1> /dev/null 2>&1;
then
       $DIRECTORY/deploy.sh $DEPLOYEDFILE 
fi
cd $DIRECTORY
   
done

echo Deployed:
cat $DEPLOYEDFILE 


