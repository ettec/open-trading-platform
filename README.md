## Open Trading Platform

An open source highly scalable platform for building cross asset execution orientated trading applications that can be easily deployed on-prem or in the cloud.

### Table of contents
1. [synopsis](#synopsis)
2. [architecture overview](#architectureoverview)
3. [services](#services)
4. [so what use is it to you?](#sowhatuseisittoyou)
5. [about the author](#abouttheauthor)

## synopsis <a name="synopsis"></a>

The platform consists of a number of services that are common in execution orientated trading applications, which will we get to shortly.  Before that I would like to outline the approach behind the platform.  If however you just want to get an out of the box configuration of the platform running see [here](https://github.com/ettec/open-trading-platform/blob/master/install/README.md) for the simple installation guide.

Historically the choice presented when a firm has needed it's own trading platform has been between choosing a vendor solution or building a bespoke solution.  The vendor platform option usually lacks flexibility and ends up costing materially more than initially allowed for once the costs of running and customising the platform are considered.  The bespoke solution is highly flexible but has a large upfront development cost and can have a higher ongoing maintenance cost (though this is not necessarily the case as platform specialists are usually required to modify and maintain the vendor solution). 

 The aim of OTP is to offer a 3rd option that gives the flexibility of a bespoke solution without the high upfront development costs and with reduced maintenance costs.  OTP is released under a GPLv3 license meaning it is and always will be free.  The platform makes extensive use of tried and tested open source projects to meet the, frankly more complex, non-functional requirements of a typical trading platform, this reduces maintenance costs and enhances the robustness of the system.  Through the extensive use of open source projects the ratio of OTP code to functionality is low which, importantly, also means the barrier to understanding and modifying the OTP code is low.

It is the case that most vendor trading platforms around today were written before the open source movement really got off the ground and thus they have a sunk cost in non-functional code and ongoing maintenance costs which would not have been the case had they had the option to leverage some of the world class open source solutions now available.  The details of the key open source projects used in OTP are outlined in the next section, but the guiding philosophy when deciding whether to use an open source project throughout the building of this platform has been much akin to that of the requirements for a [CNCF](https://www.cncf.io/) graduated project (of which a number are used in OTP), which in essence can be summarised as:  the project is in use and maintained by multiple significant technology organisations and has been shown to follow best practice.  



## architecture overview <a name="architectureoverview"></a>

Below is the list of key technologies used in OTP and a brief but by no means comprehensive outline of what they bring to the party:

#### language overview

All server side components are written in Golang, apart from the FIX Market Simulator which is written in Java (historical reasons, i.e. it happened to be the language I knew best when I started the project).  Golang is IMO an excellent language that just gets out the way and allows the developer to focus on the problem being solved (it 'minimises cognitive friction') and I will say with confidence that it is a very good match for this sort of application.

The client is an SPA web application written in Typescript using the React library.  I have to say I have been pleasantly surprised at how productive these technologies have been to use (what little client side development experience I have with which to compare has been primarily in Java RCP and C# wpf) .

#### open source projects

**Kubernetes:** The backbone of the platform, makes it cinch to scale and deploy the platform components.  Also means it's as easy to deploy OTP in the cloud as it is to deploy it on-prem. 

**Helm:** Automates the install of platform components.

**Kafka:** Essentially a distributed transaction log, this is used to distribute order and execution information across the platform.  Using this as the backbone of the system makes scaling the order store very straightforward (through increasing the number of order topic partitions to increase I/O parallelism)

**Protobuf:** Used to define the business model and service apis of the platform and makes it easy to share the business model and service apis across both the server and client.

**gRPC:**  a polyglot compact binary communication protocol that provides a standardised way for the platform services to communicate via streams or Rpc calls.

**Envoy/grpc-web:**  the client's gateway to the OTP platform's services.  It allows both the protobuf business model and grpc service apis to be seemlessly shared by the web client.  Also supports streaming data to the client. 

**Prometheus:**  used to capture significant real-time and historical per service performance stats such as quote fan in/fan out, order counts etc. 

**Postgresql:** used as a store for static data (instruments, markets, listings) and client configuration data

**Grafana:**  a number of OTP specific dashboards are provided as part of the platform to assist with monitoring.

**BlueprintJS:**  a library of GUI components designed for applications built to analyse data; there is a good affinity between this type of application and a trading application, I have found it to be a good match in practice.

**Caplin Flexlayout:** the ideal GUI layout manager for a trading application as it was originally written to be, well, a layout manager for a trading application (thanks to Caplin for open sourcing it).



## services  <a name="services"></a>

Below is a list of links to the  source route and README file of each platform service that further elaborates upon the details of the service.  Note, a basic familiarity with Kubernetes will be beneficial to be able to understand some the terminology .

[authorization-service](https://github.com/ettec/open-trading-platform/blob/master/go/authorization-service)

[client-config-service](https://github.com/ettec/open-trading-platform/blob/master/go/client-config-service)

[fix-market-simulator](https://github.com/ettec/open-trading-platform/blob/master/java/fixmarketsimulator)

[fix-sim-execution-venue](https://github.com/ettec/open-trading-platform/blob/master/go/execution-venues/fix-sim-execution-venue)

[order-router](https://github.com/ettec/open-trading-platform/blob/master/go/execution-venues/order-router)

[smart-router](https://github.com/ettec/open-trading-platform/tree/master/go/execution-venues/smart-router)

[vwap-strategy](https://github.com/ettec/open-trading-platform/blob/master/go/execution-venues/vwap-strategy)

[market-data-gateway-fixsim](https://github.com/ettec/open-trading-platform/blob/master/go/market-data/market-data-gateway-fixsim)

[market-data-service](https://github.com/ettec/open-trading-platform/blob/master/go/market-data/market-data-service)

[quote-aggregator](https://github.com/ettec/open-trading-platform/tree/master/go/market-data/quote-aggregator)

[opentp-client](https://github.com/ettec/open-trading-platform/blob/master/react/opentp-client)

[order-data-service](https://github.com/ettec/open-trading-platform/blob/master/go/order-data-service)

[order-monitor](https://github.com/ettec/open-trading-platform/blob/master/go/order-monitor)

[static-data-service](https://github.com/ettec/open-trading-platform/blob/master/go/static-data-service)

## so what use is it to you?  <a name="sowhatuseisittoyou"></a>

If you are considering writing a bespoke trading platform OTP could be used as a starting point to give you a significant leg up.  Or perhaps the firm you work for is of a size for which a bespoke platform would not normally make sense but you would like to avoid all the costs, hassles and lack of flexibility involved in using a vendor platform, OTP may make having your own bespoke platform feasible.  Or maybe you are after some ideas and guidance on technologies to use in your own solution.  OTP, I believe, will bring value to all 3 of these situations.



## about the author <a name="abouttheauthor"></a>

Software is in my blood, I've been coding since I was 10, first in BASIC, then assembler, I took a detour via a Physics degree, but the pull of coding was just too much for me to resist and I ended up as a software engineer at a number of mainly financial organisations working primarily on front office trading systems.  For the last few years I have opted to be a stay at home dad to support my wife during her PhD and to allow me to spend more time with my kids whilst they were still young, amongst other projects.  I started tentatively experimenting with what would become OTP towards the end of 2019 initially just as a POC and way of exploring technologies, however the more I built the more I could see that the unique combination of technologies could yield a genuinely innovative trading platform.  If you are considering using all or part of it (it's GPLv3 licensed so is and always will be free) I'd be happy to help you get off the ground.  Contact me at matthew.pendrey@gmail.com   

