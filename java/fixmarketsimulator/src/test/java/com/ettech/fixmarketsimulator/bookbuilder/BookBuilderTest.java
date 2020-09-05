package com.ettech.fixmarketsimulator.bookbuilder;

import com.google.gson.Gson;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.math.RoundingMode;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class BookBuilderTest {

    String depthJson = "[{\"symbol\":\"CHK\",\"marketPercent\":0.0292,\"volume\":3743791,\"lastSalePrice\":0.157,\"last"+
            "SaleSize\":100,\"lastSaleTime\":1583859385891,\"lastUpdated\":1583859389226,\"bids\":[{\"price\":0.157,\"s"+"" +
            "ize\":800,\"timestamp\":1583859389226},{\"price\":0.15,\"size\":300,\"timestamp\":1583858752568},{\"pr"+
            "ice\":0.13,\"size\":100,\"timestamp\":1583847222423},{\"price\":0.11,\"size\":300,\"timestamp\":1583849" +
            "622116},{\"price\":0.1,\"size\":100,\"timestamp\":1583847222420},{\"price\":0.08,\"size\":100,\"tim" +
            "estamp\":1583847222421}],\"asks\":[{\"price\":0.27,\"size\":100,\"timestamp\":1583847222416},{\"price\":0.3" +
            "8,\"size\":200,\"timestamp\":1583849255760},{\"price\":0.4,\"size\":100,\"timestamp\":1583849292063},{\"p" +
            "rice\":0.42,\"size\":100,\"timestamp\":1583847222415}],\"systemEvent\":{\"systemEvent\":\"R\",\"times" +
            "tamp\":1583847000000},\"tradingStatus\":{\"status\":\"T\",\"reason\":\"    \",\"timestamp\":15838401338" +
            "12},\"opHaltStatus\":{\"isHalted\":false,\"timestamp\":1583840133812},\"ssrStatus\":{\"isSSR\":true,\"det" +
            "ail\":\"N\",\"timestamp\":1583840348320},\"securityEvent\":{\"securityEvent\":\"MarketOpen\",\"timest" +
            "amp\":1583847000000},\"trades\":[],\"tradeBreaks\":[]},{\"symbol\":\"XLF\",\"marketPercent\":0.04207,\"v" +
            "olume\":2652034,\"lastSalePrice\":23.215,\"lastSaleSize\":600,\"lastSaleTime\":1583859406852,\"lastU" +
            "pdated\":1583859408939,\"bids\":[{\"price\":23.21,\"size\":600,\"timestamp\":1583859407049},{\"price\":23" +
            ".19,\"size\":2600,\"timestamp\":1583859395793},{\"price\":23.18,\"size\":2500,\"timestamp\":15838" +
            "59400210},{\"price\":23.16,\"size\":100,\"timestamp\":1583859395672},{\"price\":22.87,\"size\":300,\"timesta" +
            "mp\":1583856001265},{\"price\":22.86,\"size\":300,\"timestamp\":1583855906095}],\"asks\":[{\"price\":23." +
            "23,\"size\":2500,\"timestamp\":1583859398871},{\"price\":23.24,\"size\":2500,\"timestamp\":1583859395814}" +
            ",{\"price\":23.25,\"size\":2500,\"timestamp\":1583859396170},{\"price\":23.47,\"size\":200,\"timestamp\":1" +
            "583852391071}],\"systemEvent\":{\"systemEvent\":\"R\",\"timestamp\":1583847000000},\"tradingStatus\":{\"s" +
            "tatus\":\"T\",\"reason\":\"    \",\"timestamp\":1583840133821},\"opHaltStatus\":{\"isHalted\":false,\"time" +
            "stamp\":1583840133821},\"ssrStatus\":{\"isSSR\":true,\"detail\":\"N\",\"timestamp\":1583840371079},\"secur" +
            "ityEvent\":{\"securityEvent\":\"MarketOpen\",\"timestamp\":1583847000000},\"trades\":[],\"tradeBreaks\":[" +
            "]},{\"symbol\":\"USO\",\"marketPercent\":0.04999,\"volume\":2561681,\"lastSalePrice\":7.1,\"lastSal" +
            "eSize\":500,\"lastSaleTime\":1583859364361,\"lastUpdated\":1583859406935,\"bids\":[{\"price\":7.09,\"si" +
            "ze\":2400,\"timestamp\":1583859406935},{\"price\":7.08,\"size\":12500,\"timestamp\":1583859405300},{\"pr" +
            "ice\":7.07,\"size\":12500,\"timestamp\":1583859381431},{\"price\":7.06,\"size\":500,\"timestamp\":158385923" +
            "3786},{\"price\":7.05,\"size\":500,\"timestamp\":1583859228553},{\"price\":7.04,\"size\":500,\"timestamp\":158" +
            "3859375257},{\"price\":6.96,\"size\":100,\"timestamp\":1583858080117},{\"price\":6.87,\"size\":145,\"times" +
            "tamp\":1583856700384},{\"price\":6.86,\"size\":1600,\"timestamp\":1583857896390},{\"price\":6.85,\"size\":40" +
            "000,\"timestamp\":1583855705649},{\"price\":6.8,\"size\":100,\"timestamp\":1583852609754},{\"price\":6.79,\"si" +
            "ze\":100,\"timestamp\":1583847002783},{\"price\":6.77,\"size\":100,\"timestamp\":1583853005853},{\"price\":6.7" +
            "3,\"size\":3000,\"timestamp\":1583847004159},{\"price\":6.7,\"size\":900,\"timestamp\":1583847016002},{\"pri" +
            "ce\":6.68,\"size\":100,\"timestamp\":1583852271221},{\"price\":6.55,\"size\":100,\"timestamp\":158384908" +
            "3119},{\"price\":6.5,\"size\":1900,\"timestamp\":1583853217368},{\"price\":6.2,\"size\":100,\"timestamp\":158384" +
            "7002769},{\"price\":6.08,\"size\":600,\"timestamp\":1583856171789},{\"price\":6,\"size\":100,\"timestamp\":158" +
            "3849128742},{\"price\":5.5,\"size\":300,\"timestamp\":1583847028584},{\"price\":5.23,\"size\":400,\"timestamp\":1" +
            "583847028367},{\"price\":5.1,\"size\":100,\"timestamp\":1583847031715},{\"price\":4.88,\"size\":600,\"timesta" +
            "mp\":1583847027704}],\"asks\":[{\"price\":7.1,\"size\":12500,\"timestamp\":1583859379471},{\"price\":7.11,\"si" +
            "ze\":12500,\"timestamp\":1583859365167},{\"price\":7.12,\"size\":500,\"timestamp\":1583859375011},{\"price\":7" +
            ".13,\"size\":500,\"timestamp\":1583859229623},{\"price\":7.14,\"size\":500,\"timestamp\":1583859248800},{\"pri" +
            "ce\":7.24,\"size\":100,\"timestamp\":1583850109805},{\"price\":7.48,\"size\":162,\"timestamp\":158385338723" +
            "5}],\"systemEvent\":{\"systemEvent\":\"R\",\"timestamp\":1583847000000},\"tradingStatus\":{\"status\":\"T\",\"" +
            "\":\"    \",\"timestamp\":1583840133820},\"opHaltStatus\":{\"isHalted\":false,\"timestamp\":1583840133820},\"ss" +
            "rStatus\":{\"isSSR\":true,\"detail\":\"N\",\"timestamp\":1583840371070},\"securityEvent\":{\"securityEvent\":\"Ma" +
            "rketOpen\",\"timestamp\":1583847000000},\"trades\":[],\"tradeBreaks\":[]}]";


    @Test
    void testDepthParsing() throws Exception {
        Gson g = new Gson();
        Depth[] depths = g.fromJson(depthJson, Depth[].class);

        assertEquals(depths.length, 3);
    }



}
