
FROM openjdk:latest

ADD target/fixmarketsimulator-1.0-SNAPSHOT-bin.tar.gz /usr/src
ADD swagger /swagger

CMD java -cp "/usr/src/fixmarketsimulator-1.0-SNAPSHOT/lib/*" com.ettech.fixmarketsimulator.App