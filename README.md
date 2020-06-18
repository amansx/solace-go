# solace-go

solclient -> solace client library
solwrap   -> solace wrapper c++ library that links with solclient
gosol     -> golang bare metal binding to solwrap
gosol.lib -> golang high level binding as plugin

## Compilation steps
1. solwrap/x.cpp + solclient.a -> libsolwrap.a
2. gosol.go + libsolwrap.a -> solace.plugin
3. application -> solace.plugin

# Solace Concepts
## [Direct Messaging](https://docs.solace.com/PubSub-Basics/Direct-Messages.htm)
Direct messaging provides a reliable, but not guaranteed, delivery of messages from the Solace message bus to consuming clients, 
and is the default message delivery system for Solace PubSub+. 
It doesn't require any configuration beyond that required to set up and start other features. 
By default, direct messaging is always available to all clients connecting to an event broker.

### Characteristics of direct messages
* Aren't retained for a client when it's not connected to an event broker.
* Delivered to subscribing clients in the order in which publishers publish them.
* Don't require acknowledgment of receipt by subscribing clients.
* Aren't spooled on the message bus for consuming clients.
* Can be discarded when congestion or system failures are encountered.

### Typical applications
* Extremely high message rates and extremely low latency are required.
* Consuming clients can tolerate message loss in the event of network congestion.
* Messages don't need to be persisted for later delivery to slow or offline consumers.
* Messages must be efficiently published to a large number of clients with matching subscriptions.

### Situations where direct messages aren't delivered
* A client disconnects from an event broker while messages are being published.
* A client fails to consume messages at the published rate and egress message buffer overflow results.


## [Guaranteed Messaging](https://docs.solace.com/PubSub-Basics/Guaranteed-Messages.htm)
Guaranteed Messaging can be used to ensure the delivery of a message between two applications even in cases where 
the receiving application is off line, or there is a failure of a piece of network equipment. 
As well, those messages are delivered in the order they were published.

### Key points about Guaranteed Messaging:
* It keeps messages that are tagged as persistent or non-persistent (rather than Direct) across event broker restarts by spooling them to persistent storage.
* It keeps a copy of the message until successful delivery to all clients and downstream event brokers has been verified.
* Messages accepted by Solace PubSub+ through Guaranteed Messaging for delivery to clients are never lost, but might not get accepted if system resource limits are exceeded. 
* If an ingress message cannot be received by Solace PubSub+ (for example, if the spool quota is exceeded), the publisher is not acknowledged, and the appropriate event broker statistic is incremented.

*Note*: To support Guaranteed Messaging, an event broker must have message spooling enabled, and a Solace PubSub+ appliance must also have an Assured Delivery Blade (ADB) installed.

*Note*: Although a persistent delivery mode is typically used for Guaranteed messages, a non-persistent delivery mode is provided to offer compatibility with Java Message Service (JMS), 
and to allow the delivery mode of the messages to be modified to accommodate the persistence requirements of an endpoint or a client subscription when there is a topic match 
(refer to Topic Matching & Message Delivery Modes).



## Running Solace Locally

### Salient Points
* The message spool is configured and enabled. The maximum spool usage is set to 1500 MB.
* The message VPN named default is enabled with no client authentication.
* The client username called default in the default Message VPN is enabled. 
* The client username default has all settings enabled.
* All features are unlocked and do not require a product key.
* All services are enabled and the table below lists their default port numbers.

### Useful Ports
* 8080: Enables SEMP management traffic to the container. Use this port when connecting to the container using the PubSub+ Manager.
* 55555: Enables SMF data to pass through the container.

### More Ports

| Service description | Type of Traffic                                                                                                                                                               |
| ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 22                  | Host OS SSH

(2222 for Solace PubSub+ software event broker Machine & Cloud Image releases prior to 8.5.0)                                                                    | Management |
| 2222                | Solace CLI SSH/SFTP

(22 for for Solace PubSub+ software event broker Machine & Cloud Image releases prior to 8.5.0)                                                          | Management |
| 8080                | Access to PubSub+ Manager, SEMP, and SolAdmin                                                                                                                                 | Management |
| 1943                | Access to PubSub+ Manager over HTTPS, SEMP over TLS, and SolAdmin over TLS.

For Solace PubSub+ software event broker versions before the 9.4.0 release, use port 943.        | Management |
| 5550                | Health Check Listen Port                                                                                                                                                      | Data |
| 55555               | Solace Message Format (SMF)                                                                                                                                                   | Data |
| 55003               | SMF compressed                                                                                                                                                                | Data |
| 55556               | SMF routing                                                                                                                                                                   | Data |
| 55443               | SMF TLS / SSL (with or without compression)                                                                                                                                   | Data |
| 8008                | Web Transport - WebSockets, Comet, etc.

For Solace PubSub+ software event broker versions before the 9.4.0 release, use port 80.                                             | Data |
| 1443                | Web Transport TLS / SSL

For Solace PubSub+ software event broker versions before the 9.4.0 release, use port 443.                                                            | Data |
| 5671                | AMQP TLS / SSL ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of AMQP ports)              | Data |
| 5672                | AMQP ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of AMQP ports)                        | Data |
| 1883                | MQTT ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports)                        | Data |
| 8883                | MQTT TLS / SSL('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports)               | Data |
| 8000                | MQTT / WebSockets ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports)           | Data |
| 8443                | MQTT / WebSockets TLS / SSL ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports) | Data |
| 9000                | REST ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of REST ports)                        | Data |
| 9443                | REST TLS / SSL ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of REST ports)              | Data |
| 8741                | High Availability (HA) Mate Link                                                                                                                                              | HA |
| 8300

8301

8302    | HA Configuration Synchronization                                                                                                                                              | HA |