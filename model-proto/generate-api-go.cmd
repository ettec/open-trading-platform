# Usage: generate-go.cmd  <servicename> 
SVC_NAME=$1
SVC_PATH=../go/$SVC_NAME
API_PATH=$SVC_PATH/api
GOFILE=$API_PATH/$SVC_NAME.pb.go
mkdir -p $API_PATH
protoc $SVC_PATH/$SVC_NAME.proto --go_out=plugins=grpc:$SVC_PATH/api/ --proto_path=$SVC_PATH:.

sed -i 's/Empty/model.Empty/g' $GOFILE
sed -i 's/ClobQuote/model.ClobQuote/g' $GOFILE
sed -i 's/*Listing,/*model.Listing,/g' $GOFILE
sed -i 's/*Listing /*model.Listing /g' $GOFILE
sed -i 's/*Listing)/*model.Listing)/g' $GOFILE
sed -i 's/Decimal64/model.Decimal64/g' $GOFILE
sed -i 's/ Side/ model.Side/g' $GOFILE
sed -i 's/import (/import (\n\t"github.com\/ettec\/open-trading-platform\/go\/model"/g' $GOFILE
