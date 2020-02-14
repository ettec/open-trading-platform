# The fix spec proto files provided by fixtrading.org do not compile in go, this script fixes the issues in the generated go files
# Usage: generate-go.cmd  <servicename> 
SVC_PATH=../go/$1



# Fix sim service
FIXSIM_PATH=$SVC_PATH/internal/connections/fixsim/
mkdir -p $FIXSIM_PATH
protoc $SVC_PATH/fixsimmarketdataservice.proto --go_out=plugins=grpc:$FIXSIM_PATH --proto_path=$SVC_PATH:.
GOFILE=$FIXSIM_PATH/fixsimmarketdataservice.pb.go
sed -i 's/MarketDataRequest/marketdata.MarketDataRequest/g' $GOFILE
sed -i 's/MarketDataIncrementalRefresh/marketdata.MarketDataIncrementalRefresh/g' $GOFILE
sed -i 's/import (/import (\n\tmarketdata \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/marketdata\"/g' $GOFILE




# Fix Protocol
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
sed -i 's/*UndInstrmtGrp/*common.UndInstrmtGrp/g' $GOFILE
sed -i 's/*InstrmtMDReqGrp/*common.InstrmtMDReqGrp/g' $GOFILE
sed -i 's/*TrdgSesGrp/*common.TrdgSesGrp/g' $GOFILE
sed -i 's/*RoutingGrp/*common.RoutingGrp/g' $GOFILE
sed -i 's/*ApplicationSequenceControl/*common.ApplicationSequenceControl/g' $GOFILE

sed -i 's/*StandardHeader/*session.StandardHeader/g' $GOFILE
sed -i 's/*StandardTrailer/*session.StandardTrailer/g' $GOFILE



sed -i 's/import (/import (\n\tfix \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/fix\"/g' $GOFILE
sed -i 's/import (/import (\n\tcommon \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/common\"/g' $GOFILE
sed -i 's/import (/import (\n\session \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/session\"/g' $GOFILE

mkdir -p $SVC_PATH/internal/fix/meta
protoc ./meta.proto --go_out=$SVC_PATH/internal/fix/meta --proto_path=$SVC_PATH:.

mkdir -p $SVC_PATH/internal/fix/session
protoc ./session.proto --go_out=$SVC_PATH/internal/fix/session --proto_path=$SVC_PATH:.
GOFILE=$SVC_PATH/internal/fix/session/session.pb.go
sed -i 's/*MsgTypeGrp/*common.MsgTypeGrp/g' $GOFILE
sed -i 's/Timestamp/fix.Timestamp/g' $GOFILE
sed -i 's/import (/import (\n\tfix \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/fix\"/g' $GOFILE
sed -i 's/import (/import (\n\tcommon \"github.com\/ettec\/open-trading-platform\/go\/market-data-gateway\/internal\/fix\/common\"/g' $GOFILE



