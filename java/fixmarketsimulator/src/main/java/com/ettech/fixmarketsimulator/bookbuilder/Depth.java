package com.ettech.fixmarketsimulator.bookbuilder;

import java.util.List;



public class Depth{
    public String symbol;
    public double marketPercent;
    public int volume;
    public double lastSalePrice;
    public int lastSaleSize;
    public long lastSaleTime;
    public long lastUpdated;
    public List<Line> bids;
    public List<Line> asks;
    public SystemEvent systemEvent;
    public TradingStatus tradingStatus;
    public OpHaltStatus opHaltStatus;
    public SsrStatus ssrStatus;
    public SecurityEvent securityEvent;
    public List<Object> trades;
    public List<Object> tradeBreaks;

    public static class Line {
        public double price;
        public int size;
        public Object timestamp;
    }

    public static  class SystemEvent{
        public String systemEvent;
        public long timestamp;
    }

    public  static class TradingStatus{
        public String status;
        public String reason;
        public long timestamp;
    }

    public static  class OpHaltStatus{
        public boolean isHalted;
        public long timestamp;
    }

    public  static class SsrStatus{
        public boolean isSSR;
        public String detail;
        public long timestamp;
    }

    public  static  class SecurityEvent{
        public String securityEvent;
        public long timestamp;
    }
}
