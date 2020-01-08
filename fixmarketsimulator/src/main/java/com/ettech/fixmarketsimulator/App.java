package com.ettech.fixmarketsimulator;

import com.ettech.fixmarketsimulator.exchange.ExchangeSimulatorGuiceModule;
import com.ettech.fixmarketsimulator.fix.ApplicationImpl;
import com.ettech.fixmarketsimulator.marketdataserver.MarketDataService;
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

    private static int fixServerPort = 9876;


    private static String fixConfig =
            "[default]\n" +
                    // "FileStorePath=/usr/share/cnoms/fixmarketsimulator\n" +
                    "FileStorePath=" + System.getenv("FIX_FILE_STORE_PATH") + "\n" +
                    //"FileStorePath=./fixlog\n" +
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


    public static void main(String[] args) {
        int port = 8501;

        log.info("Starting api server on port:" + port);
        log.info("Starting fix server on port:" + fixServerPort);

        exchangeSimulatorInjector = Guice.createInjector(new ExchangeSimulatorGuiceModule());


        Server server = new Server(port);

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


        Acceptor acceptor = null;
        try {

            Application application = exchangeSimulatorInjector.getInstance(ApplicationImpl.class);

            InputStream stream = new ByteArrayInputStream(fixConfig.getBytes(StandardCharsets.UTF_8));

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
}
