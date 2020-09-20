package com.ettech.fixmarketsimulator;

import com.ettech.fixmarketsimulator.bookbuilder.BooksBuilder;
import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.ExchangeSimulatorGuiceModule;
import com.ettech.fixmarketsimulator.exchange.impl.ExchangeImpl;
import com.ettech.fixmarketsimulator.fix.ApplicationImpl;
import com.ettech.fixmarketsimulator.marketdataserver.MarketDataService;
import com.ettech.fixmarketsimulator.orderentryserver.OrderEntryService;
import com.google.inject.Guice;
import com.google.inject.Injector;
import io.swagger.v3.jaxrs2.integration.JaxrsOpenApiContextBuilder;
import io.swagger.v3.oas.integration.OpenApiConfigurationException;
import io.swagger.v3.oas.integration.SwaggerConfiguration;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import org.eclipse.jetty.server.Handler;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.server.handler.DefaultHandler;
import org.eclipse.jetty.server.handler.HandlerList;
import org.eclipse.jetty.server.handler.ResourceHandler;
import org.eclipse.jetty.servlet.ServletContextHandler;
import org.eclipse.jetty.servlet.ServletHolder;
import org.glassfish.jersey.servlet.ServletContainer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import quickfix.*;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class App {
    static Logger log = LoggerFactory.getLogger(App.class);

    public static Injector exchangeSimulatorInjector;

    static MarketDataService marketDataServer;
    static OrderEntryService orderEntryServer;
    static BooksBuilder booksBuilder;
    static Exchange exchange;

    private static int fixServerPort;

    static {
        fixServerPort = Integer.parseInt(getSysEnvVal("FIX_SERVER_PORT", "9876"));
    }


    private static String getFixConfig() {
        var config = "[default]\n" +
                "FileStorePath=" + System.getenv("FIX_FILE_STORE_PATH") + "\n" +
                "SocketAcceptPort=" + fixServerPort + "\n" +
                "BeginString= FIXT.1.1\n" +
                "DefaultApplVerID= FIX.5.0SP2\n" +
                "\n";

        var targetCompIds = System.getenv("TARGET_COMP_IDS").split(",");


        for (String targetCompId : targetCompIds) {
            config +=
                    "[session]\n" +
                            "SenderCompID=EXEC\n" +
                            "TargetCompID=" + targetCompId + "\n" +
                            "ConnectionType=acceptor\n" +
                            "StartTime=00:00:00\n" +
                            "FileLogPath=fixlog\n" +
                            "EndTime=00:00:00\n";

        }

        return config;
    }




    /*
    private static String fixConfig =
            "[default]\n" +
                    "FileStorePath=" + System.getenv("FIX_FILE_STORE_PATH") + "\n" +
                    "SocketAcceptPort=" + fixServerPort + "\n" +
                    "BeginString= FIXT.1.1\n" +
                    "DefaultApplVerID= FIX.5.0SP2\n" +
                    "\n" +
                    "[session]\n" +
                    "SenderCompID=EXEC\n" +
                    "TargetCompID=BANZAI\n" +
                    "ConnectionType=acceptor\n" +
                    "StartTime=00:00:00\n" +
                    "FileLogPath=fixlog\n" +
                    "EndTime=00:00:00";
*/

    public static void main(String[] args) {
        int restApiPort = Integer.parseInt(getSysEnvVal("REST_API_PORT", "8501"));
        ;

        log.info("Starting rest server on restApiPort:" + restApiPort);
        log.info("Starting fix server on restApiPort:" + fixServerPort);

        exchange = new ExchangeImpl();

        exchangeSimulatorInjector = Guice.createInjector(new ExchangeSimulatorGuiceModule(exchange));

        Server server = new Server(restApiPort);

        ResourceHandler resource_handler = new ResourceHandler();
        resource_handler.setDirectoriesListed(true);
        resource_handler.setWelcomeFiles(new String[]{"index.html"});

        resource_handler.setResourceBase("./swagger");

        HandlerList handlers = new HandlerList();


        ServletContextHandler ctx =
                new ServletContextHandler(ServletContextHandler.NO_SESSIONS);

        ctx.setContextPath("/");


        handlers.setHandlers(new Handler[]{resource_handler, ctx, new DefaultHandler()});
        server.setHandler(handlers);


        ServletHolder apiServlet = ctx.addServlet(ServletContainer.class, "/*");

        apiServlet.setInitOrder(1);
        apiServlet.setInitParameter("jersey.config.server.provider.packages", "com.ettech.fixmarketsimulator.api,io.swagger.v3.jaxrs2.integration.resources");


        // Setup API resources
        OpenAPI oas = new OpenAPI();
        Info info = new Info()
                .title(App.class.getPackageName());

        oas.info(info);
        SwaggerConfiguration oasConfig = new SwaggerConfiguration()
                .openAPI(oas)
                .prettyPrint(true)
                .resourcePackages(Stream.of("com.ettech.fixmarketsimulator.api").collect(Collectors.toSet()));

        try {
            JaxrsOpenApiContextBuilder builder = new JaxrsOpenApiContextBuilder();
            builder.openApiConfiguration(oasConfig).buildContext(true);
        } catch (OpenApiConfigurationException e) {
            throw new RuntimeException(e.getMessage(), e);
        }

        marketDataServer = exchangeSimulatorInjector.getInstance(MarketDataService.class);
        try {
            marketDataServer.start();
        } catch (IOException e) {
            log.error("Failed to start market data server", e);
        }

        orderEntryServer = exchangeSimulatorInjector.getInstance(OrderEntryService.class);
        try {
            orderEntryServer.start();
        } catch (IOException e) {
            log.error("Failed to start orderEntryServer", e);
        }

        try {

            booksBuilder = new BooksBuilder(exchange, getSysEnvVal("BB_DEPTH_PATH", "./resource/depth.json"),
                    getSysEnvVal("BB_SYMS_TO_RUN", "MSFT,SPY"),
                    Integer.parseInt(getSysEnvVal("BB_NUM_EXEC_THREADS", "5")),
                    Integer.parseInt(getSysEnvVal("BB_UPDATE_INTERVAL_MS", "1000")),
                    Double.parseDouble(getSysEnvVal("BB_MIN_QTY", "0.9")),
                    Double.parseDouble(getSysEnvVal("BB_VARIATION", "0.005")),
                    Integer.parseInt(getSysEnvVal("BB_TICK_SCALE", "2")),
                    Double.parseDouble(getSysEnvVal("BB_TRADE_PROBABILITY", "0.2")),
                    Integer.parseInt(getSysEnvVal("BB_MAX_DEPTH", "10")),
                    Double.parseDouble(getSysEnvVal("BB_CANCEL_PROBABILITY", "0.2")));
        } catch (Exception e) {
            throw new RuntimeException("failed to create book builder", e);
        }

        Acceptor acceptor = null;
        try {

            Application application = exchangeSimulatorInjector.getInstance(ApplicationImpl.class);

            InputStream stream = new ByteArrayInputStream(getFixConfig().getBytes(StandardCharsets.UTF_8));

            SessionSettings settings = new SessionSettings(stream);
            MessageStoreFactory storeFactory = new FileStoreFactory(settings);
            LogFactory logFactory = new FileLogFactory(settings);
            MessageFactory messageFactory = new DefaultMessageFactory();
            acceptor = new SocketAcceptor(application, storeFactory, settings, logFactory, messageFactory);


            acceptor.start();


            server.start();
            server.join();
        } catch (Exception ex) {
            log.error("Failed to start jetty server", ex);
        } finally {
            try {
                if (marketDataServer != null) {
                    marketDataServer.stop();
                    marketDataServer.blockUntilShutdown();
                }
            } catch (Throwable t) {
                log.error("error whilst closing market data server", t);
            }

            server.destroy();
        }
    }

    public static String getSysEnvVal(String key, String defaultVal) {
        var val = System.getenv(key);
        if (val == null) {
            return defaultVal;
        }

        return val;
    }
}
