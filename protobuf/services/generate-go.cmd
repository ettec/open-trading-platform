# Usage: generate-go.cmd  <path> <proto-file-name>
OUTPATH=$1
PROTO_FILE_NAME=$2


if test -z "$PROTO_FILE_NAME" 
then
      echo "second argument must be the proto file name"
      exit 1
fi
GOFILE=$OUTPATH/${PROTO_FILE_NAME/.proto/.pb.go}
protoc $PROTO_FILE_NAME --go_out=plugins=grpc:$OUTPATH --proto_path=.:../model

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
sed -i 's/import (/import (\n\t"github.com\/ettec\/otp-common\/model"/g' $GOFILE
