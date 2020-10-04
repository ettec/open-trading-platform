# market-data-service

This service implements the [market data service api](https://github.com/ettec/open-trading-platform/blob/master/protobuf/services/marketdataservice.proto).  The market data service load balances quote subscriptions across market data gateways by listing id for a given market and fans out market data to clients.  Internally it has a per client conflated queue to ensure that slow clients always get the latest quote.  The service can be scaled by increasing the deployments replica count.

