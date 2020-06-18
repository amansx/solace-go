# solace-go

| Component | Description|
|-----------|------------|
| solclient | solace client library |
| solwrap   | solace wrapper c++ library that links with solclient |
| gosol     | golang bare metal binding to solwrap |
| gosol.lib | golang high level binding as plugin |

## Compilation steps
1. solwrap/x.cpp + solclient.a -> libsolwrap.a
2. gosol.go + libsolwrap.a -> solace.plugin
3. application -> solace.plugin

# Solace Concepts
## Message Format

### [Structure](https://solace.com/blog/inside-solace-message-introduction/)
* Meta data -> Header (destination, delivery mode, return address) + Properties (application properties)
* Body -> Payload/Attachment/Application Data

### [Header and Properties](https://solace.com/blog/inside-a-solace-message-using-header-properties/)
* Binary content—Binary data can be added to a message as a binary attachment. A message can only contain a single attachment.
* When this attachment is sent through the event broker, it is not processed, transformed, or considered in subscription matching or filtering operations. This provides an efficient means for sending data that does not require processing by the platform before it reaches receiving applications.
* Structured data can also be added as a payload in the binary attachment (refer to Using Structured Data).
* User Property Map—Structured data can be added to user-defined message header fields.
* User Data—Up to 36 bytes of application‑specific binary data, known as user data, can be added to the User Data message header field.
* [Best Practices](https://docs.solace.com/Solace-PubSub-Messaging-APIs/API-Developer-Guide/C-API-Best-Practices.htm)

### [Payload](https://solace.com/blog/inside-a-solace-message-part-3-payload-data/)
* If you are sending a few headers that describe large content, consider setting the headers in the USER_PROPERTY map (set through solClient_msg_createUserPropertyMap(...)) and add the content using solClient_msg_setBinaryAttachmentPtr().
* When building a container, always try to accurately estimate the required size. The container could be a user property map (created through createUserPropertyMap()), a map (created in a message through solClient_msg_createBinaryAttachmentMap()) or a stream (created in a message through solClient_msg_createBinaryAttachmentStream(...)).
* When building a complex container that uses a submap or substream, write the submap or substream completely and call solClient_container_createStream(...) to finish the submap or substream before adding more to the main container.
* [Add Data Payload](https://docs.solace.com/Solace-PubSub-Messaging-APIs/API-Developer-Guide/Adding-Data-Payloads.htm)


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


## [Topic vs Topic Endpoint vs Queue](https://solace.com/blog/queues-vs-topic-endpoints/)
Solace endpoints are objects created on the event broker to persist messages. There are two types of endpoints: a queue endpoint (usually just called a queue) and a topic endpoint

### Topic
* Topic endpoint is not the same as a topic. 
* Topics are a message property the event broker uses to route messages to their destination, and they aren’t administratively configured on the event broker.

### Topic Endpoint
* Topic endpoints, are objects that define the storage of messages for a consuming application, and they do need to be provisioned on the event broker. 
* Topic endpoints are more closely related to queues than to topics.

### Topic Endpoint vs Queue
* **Sending by Name**: A producing application has the option to send a message directly to a queue by referencing that queue by its name in the message properties. 
* A producing application cannot, however, reference topic endpoints by name, and therefore only persist messages routed to the topic subscription applied to the topic endpoint.
* **Support for Multiple Topic Subscriptions**: Queues support multiple topic subscriptions. This allows for topic aggregation to a single consuming application. 
* Topic endpoints support a single subscription; should that subscription change, all messages persisted to the topic endpoint are deleted. Messages persisted in queues are unaffected by subscription changes.
* **Application of Selectors**: Selectors can be applied to either type of endpoint to apply conditional logic, but the two endpoints differ in when the conditional is processed. With queues, selectors are processed at egress, and with topic endpoints they are processed on ingress. This means that queues persist every message even if they don’t match the selector, 
* while with topic endpoints messages are only persisted if they match both the topic subscription and the selector.
* **Support for Multiple Consumers**: Queues support multiple consumers, providing a fault tolerant option should a consuming application disconnect from the queue. 
* Exclusive topic endpoints support a single consumer, matching the functionality of a JMS durable subscription. 
* Non-exclusive topic endpoints can support multiple consumers for load balancing purposes.
* **Ability to Read Without Removal**: Finally, a queue can be read from without removing messages, 
* whereas topic endpoints require the removal of messages to be read.

## [Durable vs Non-Durable](https://solace.com/blog/solace-endpoints-durable-vs-non-durable/)

### Non-durable Endpoints
A non-durable endpoint, also known as a temporary endpoint, has a shorter lifecycle and can only be created by an application. The endpoint has a lifespan of the client that created it, with an additional 60 seconds in case of unexpected disconnect. The 60 seconds provides the client with some time to reconnect to the endpoint before it and its contents are deleted from the Solace broker.

### Durable Endpoints
* A durable endpoint is one that stays on the Solace broker until it is explicitly deleted. 
* This endpoint can be created by an administrator, or dynamically created by an application. 
* Durable endpoints are most often used when consuming applications require all messages including those received by the event broker when the application is disconnected.
* Through the [client-profile](https://docs.solace.com/Configuring-and-Managing/Configuring-Client-Profiles.htm#Config-Guaranteed-Params) object, an administrator can separately control whether an application is able to dynamically create its own durable or non-durable endpoints, and how many it is able to create.
* Endpoints also have their own permissions to control which applications can consume from them, modify them, and remove them. 
* These permissions are set when the endpoint is initially created. 
* If a client application creates an endpoint, it is automatically the owner of that endpoint and has full access to it.


## Queues
* [Queue Access Types](https://solace.com/blog/solace-message-queue-access-types/)


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

| Service description | Type of Traffic                                                                                                                                                               | use case   |
| ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- |
| 22                  | Host OS SSH (2222 for Solace PubSub+ software event broker Machine & Cloud Image releases prior to 8.5.0)                                                                     | Management |
| 2222                | Solace CLI SSH/SFTP(22 for for Solace PubSub+ software event broker Machine & Cloud Image releases prior to 8.5.0)                                                            | Management |
| 8080                | Access to PubSub+ Manager, SEMP, and SolAdmin                                                                                                                                 | Management |
| 1943                | Access to PubSub+ Manager over HTTPS, SEMP over TLS, and SolAdmin over TLS.  For Solace PubSub+ software event broker versions before the 9.4.0 release, use port 943.        | Management |
| 5550                | Health Check Listen Port                                                                                                                                                      | Data |
| 55555               | Solace Message Format (SMF)                                                                                                                                                   | Data |
| 55003               | SMF compressed                                                                                                                                                                | Data |
| 55556               | SMF routing                                                                                                                                                                   | Data |
| 55443               | SMF TLS / SSL (with or without compression)                                                                                                                                   | Data |
| 8008                | Web Transport - WebSockets, Comet, etc.  For Solace PubSub+ software event broker versions before the 9.4.0 release, use port 80.                                             | Data |
| 1443                | Web Transport TLS / SSL For Solace PubSub+ software event broker versions before the 9.4.0 release, use port 443.                                                             | Data |
| 5671                | AMQP TLS / SSL ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of AMQP ports)              | Data |
| 5672                | AMQP ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of AMQP ports)                        | Data |
| 1883                | MQTT ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports)                        | Data |
| 8883                | MQTT TLS / SSL('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports)               | Data |
| 8000                | MQTT / WebSockets ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports)           | Data |
| 8443                | MQTT / WebSockets TLS / SSL ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of MQTT ports) | Data |
| 9000                | REST ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of REST ports)                        | Data |
| 9443                | REST TLS / SSL ('default' VPN; note that each message VPN configured on the Solace PubSub+ software event broker would require its own unique set of REST ports)              | Data |
| 8741                | High Availability (HA) Mate Link                                                                                                                                              | HA |
| 8300, 8301, 8302    | HA Configuration Synchronization                                                                                                                                              | HA |


## Unknowns
* Additional Data in Meta? What data?
* Who creates Queue, Default Queue Name?
* Acknowledgement mode - auto vs manual?
* Streaming?