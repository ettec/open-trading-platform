# Usage: generate-go.cmd  <servicename to generate code into> <proto-file-name without proto file extention>
SVC_NAME=$1
PROTO_FILE_NAME=$2
GO_PATH=../go
SVC_PATH=$GO_PATH/$SVC_NAME

if test -z "$PROTO_FILE_NAME" 
then
      echo "second argument must be the proto file name"
      exit 1
fi

PROTO_PATH=$GO_PATH/proto
PROTO_FILE_PATH=$PROTO_PATH/$PROTO_FILE_NAME.proto	
      
echo "using proto file: $PROTO_FILE_PATH"

API_PATH=$SVC_PATH/api/$PROTO_FILE_NAME
GOFILE=$API_PATH/$PROTO_FILE_NAME.pb.go
mkdir -p $API_PATH
protoc $PROTO_FILE_PATH --go_out=plugins=grpc:$API_PATH --proto_path=$PROTO_PATH:.

sed -i 's/Empty/model.Empty/g' $GOFILE
sed -i 's/ClobQuote/model.ClobQuote/g' $GOFILE
sed -i 's/*Listing,/*model.Listing,/g' $GOFILE
sed -i 's/(Listing)/(model.Listing)/g' $GOFILE
sed -i 's/*Listing /*model.Listing /g' $GOFILE
sed -i 's/*Listing)/*model.Listing)/g' $GOFILE
sed -i 's/Order)/model.Order)/g' $GOFILE
sed -i 's/*Order,/*model.Order,/g' $GOFILE
sed -i 's/*Order /*model.Order /g' $GOFILE
sed -i 's/*Order)/*model.Order)/g' $GOFILE
sed -i 's/Decimal64/model.Decimal64/g' $GOFILE
sed -i 's/ Side/ model.Side/g' $GOFILE
sed -i 's/*Timestamp/*model.Timestamp/g' $GOFILE
sed -i 's/import (/import (\n\t"github.com\/ettec\/open-trading-platform\/go\/model"/g' $GOFILE
