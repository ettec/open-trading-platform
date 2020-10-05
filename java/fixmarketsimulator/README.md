# fixmarketsimulator

Simulates a central limit order book.  Exposes a FIX api for order entry and FIX MD api over gRpc for market data.  In addition a swagger api is provided to submit orders and view order book state.  The simulator can be configured on a per instrument basis to simulate a live trading book against which to trade and act as a source of market data.  