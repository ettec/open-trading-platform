package com.ettech.fixmarketsimulator.exchange;

import com.ettech.fixmarketsimulator.exchange.impl.ExchangeImpl;
import com.ettech.fixmarketsimulator.marketdata.MarketDataServer;
import com.google.inject.AbstractModule;

public class ExchangeSimulatorGuiceModule extends AbstractModule {

    @Override
    protected void configure() {

        bind(Exchange.class).to( ExchangeImpl.class).asEagerSingleton();

    }
}
