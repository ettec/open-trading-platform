package com.ettech.fixmarketsimulator.marketdataserver;

import com.ettech.fixmarketsimulator.exchange.MDEntry;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class MarketDataSubscriptionTest {

    @Test
    void bigDecimalToFixDecimalConversion() {
        var fd = MarketDataSubscription.getFixDecimal64(new BigDecimal("19.36"));

        assertEquals(fd.getMantissa(), 1936);
        assertEquals(fd.getExponent(), -2);


         fd = MarketDataSubscription.getFixDecimal64(new BigDecimal("-19.36"));

        assertEquals(fd.getMantissa(), -1936);
        assertEquals(fd.getExponent(), -2);

        fd = MarketDataSubscription.getFixDecimal64(new BigDecimal("0.00004"));

        assertEquals(fd.getMantissa(), 4);
        assertEquals(fd.getExponent(), -5);

        fd = MarketDataSubscription.getFixDecimal64(new BigDecimal("1343522.2104"));

        assertEquals(fd.getMantissa(), 13435222104l);
        assertEquals(fd.getExponent(), -4);

        fd = MarketDataSubscription.getFixDecimal64(new BigDecimal("-31343522.21045432"));

        assertEquals(fd.getMantissa(), -3134352221045432l);
        assertEquals(fd.getExponent(), -8);


        fd = MarketDataSubscription.getFixDecimal64(new BigDecimal("13435222104"));

        assertEquals(fd.getMantissa(), 13435222104l);
        assertEquals(fd.getExponent(), 0);

        fd = MarketDataSubscription.getFixDecimal64(new BigDecimal("0"));

        assertEquals(fd.getMantissa(), 0l);
        assertEquals(fd.getExponent(), 0);

    }



}