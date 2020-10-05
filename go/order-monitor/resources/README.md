# order-monitor

This service implements the [order monitor api](https://github.com/ettec/open-trading-platform/blob/master/protobuf/services/ordermonitor.proto).  The order monitor tracks all platform order updates and publishes summary statistics to prometheus.  In addition it provides an api to cancel all orders for a given originator such as a trading desk or trading strategy.  

