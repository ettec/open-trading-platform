FROM maven:3.6.3-openjdk-11 AS MAVEN_BUILD

COPY ./ ./

RUN mvn clean
RUN mvn assembly::assembly

FROM openjdk:14-jdk


ADD swagger /swagger

COPY --from=MAVEN_BUILD target/fixmarketsimulator-1.0-bin.tar.gz /usr/src
RUN cd /usr/src && tar -xzf fixmarketsimulator-1.0-bin.tar.gz && rm fixmarketsimulator-1.0-bin.tar.gz

CMD java -cp "/usr/src/fixmarketsimulator-1.0/lib/*" com.ettech.fixmarketsimulator.App
