# Usage: generate-go.cmd  <servicename> 
SVC_PATH=../go/$1

mkdir -p $SVC_PATH/internal/grpc/
protoc $SVC_PATH/*.proto --go_out=plugins=grpc:$SVC_PATH/internal/grpc/ --proto_path=$SVC_PATH:.

mkdir -p $SVC_PATH/internal/fix/fix
protoc ./fix.proto --go_out=$SVC_PATH/internal/fix/fix --proto_path=$SVC_PATH:.




mkdir -p $SVC_PATH/internal/fix/common
GOFILE=$SVC_PATH/internal/fix/common/common.pb.go
protoc ./common.proto --go_out=$SVC_PATH/internal/fix/common --proto_path=$SVC_PATH:.
sed -i 's/Decimal64/fix.Decimal64/g' $GOFILE
sed -i 's/Timestamp/fix.Timestamp/g' $GOFILE
sed -i 's/ Tenor/ fix.Tenor/g' $GOFILE
sed -i 's/TrdRegfix./TrdRegfix_/g' $GOFILE
sed -i 's/LocalTimeOnly/fix.LocalTimeOnly/g' $GOFILE
sed -i 's/ Tenor/fix.Tenor/g' $GOFILE
sed -i 's/*Tenor/*fix.Tenor/g' $GOFILE
sed -i 's/*TimeOnly/*fix.TimeOnly/g' $GOFILE
sed -i 's/import (/import (\n\tfix \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/fix\"/g' $GOFILE 



mkdir -p $SVC_PATH/internal/fix/marketdata
GOFILE=$SVC_PATH/internal/fix/marketdata/marketdata.pb.go
protoc ./marketdata.proto --go_out=$SVC_PATH/internal/fix/marketdata --proto_path=$SVC_PATH:.
sed -i 's/Decimal64/fix.Decimal64/g' $GOFILE
sed -i 's/Timestamp/fix.Timestamp/g' $GOFILE
sed -i 's/ Tenor/ fix.Tenor/g' $GOFILE
sed -i 's/TrdRegfix./TrdRegfix_/g' $GOFILE
sed -i 's/LocalTimeOnly/fix.LocalTimeOnly/g' $GOFILE
sed -i 's/ Tenor/fix.Tenor/g' $GOFILE
sed -i 's/*Tenor/*fix.Tenor/g' $GOFILE
sed -i 's/*TimeOnly/*fix.TimeOnly/g' $GOFILE

sed -i 's/*SpreadOrBenchmarkCurveData/*common.SpreadOrBenchmarkCurveData/g' $GOFILE
sed -i 's/*YieldData/*common.YieldData/g' $GOFILE
sed -i 's/*RateSource/*common.RateSource/g' $GOFILE
sed -i 's/*Parties/*common.Parties/g' $GOFILE
sed -i 's/*InstrmtLegGrp/*common.InstrmtLegGrp/g' $GOFILE
sed -i 's/*Instrument/*common.Instrument/g' $GOFILE


sed -i 's/import (/import (\n\tfix \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/fix\"/g' $GOFILE
sed -i 's/import (/import (\n\tfix \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/common\/common\"/g' $GOFILE

mkdir -p $SVC_PATH/internal/fix/meta
protoc ./meta.proto --go_out=$SVC_PATH/internal/fix/meta --proto_path=$SVC_PATH:.

mkdir -p $SVC_PATH/internal/fix/session
protoc ./session.proto --go_out=$SVC_PATH/internal/fix/session --proto_path=$SVC_PATH:.
GOFILE=$SVC_PATH/internal/fix/session/session.pb.go
sed -i 's/MsgTypeGrp/common.MsgTypeGrp/g' $GOFILE
sed -i 's/Timestamp/fix.Timestamp/g' $GOFILE
sed -i 's/import (/import (\n\tfix \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/fix\"/g' $GOFILE
sed -i 's/import (/import (\n\tfix \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/common\/common\"/g' $GOFILE





#fix "github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/fix"

