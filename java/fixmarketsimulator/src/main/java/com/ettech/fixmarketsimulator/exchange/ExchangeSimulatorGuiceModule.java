package com.ettech.fixmarketsimulator.exchange;

import com.ettech.fixmarketsimulator.exchange.impl.ExchangeImpl;
import com.google.inject.AbstractModule;

public class ExchangeSimulatorGuiceModule extends AbstractModule {

    Exchange exchange;
    public ExchangeSimulatorGuiceModule(Exchange exchange) {
        this.exchange = exchange;
    }

    @Override
    protected void configure() {
        bind(Exchange.class).toInstance( this.exchange );
    }
}
