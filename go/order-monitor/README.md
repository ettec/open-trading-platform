# order-monitor

The order monitor tracks all platform order updates and publishes summary statistics to prometheus (grafana dashboards to monitor these statistics can be found [here](https://github.com/ettec/open-trading-platform/tree/master/grafana-dashboards)).  In addition it provides an api that can be used to cancel all orders for a given originator, for example a trading desk or trading strategy.

