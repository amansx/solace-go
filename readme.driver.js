/* 
Package gosol provides acces to the solace messaging platform for 
pub/sub,
queueing and 
guaranteed messaging capabilities.

All message delivery to message-consumer applications is asynchronous,
via a registered MsgHandler event callback. 

The entire package is mainly an asynchronous API, providing an interface by which applications can
register several different asynchronous event-handlers for notifications of 
Message, 
Error, 
Connectivity and 
Publisher events. 

These events are registered with the underlying Solace messaging API 
by passing function-references into the Init() function, 
which should always be the first function called in any Solace messaging client session.

The API supports several common messaging use cases such as:

- Publish/Subscribe messaging that decouples producers and consumers from each other, 
usually over a non-guaranteed transport but also supported for Guaranteed transport

- Point-to-point queueing of guaranteed messages via a named queue known 
by the publisher and subscriber applications

- Streaming publication of guaranteed messages with multiple messages 
in-flight that can be correlated back to the original message via asynchronous 
message-acknowledgements to the publisher

- Guaranteed message consumption via durable subscribers acknowledging 
messages to the Solace message broker either automatically within the 
API or manually by the application

- Guaranteed message consumption via either Solace forwarding mode: 
STORE-AND-FORWARD, or CUT-THROUGH-PERSISTENCE

*/



// Session handle is return from Init(), and passed into all other functions




/*
MsgHandler functions are registered callbacks invoked upon receipt
of Solace messages for all transports (Guaranteed and Direct). See also
type MsgEvent.
*/


/*
ErrHandler functions are registered callbacks invoked upon Solace
errors. See also type ErrEvent.
*/


/*
PubHandler functions are registered callbacks invoked upon Solace
Publisher events such as message ACK or REJECT events; for streaming
publishers it can include publisher-provided correlation data.. See also
type PubEvent.
*/


/*
ConHandler functions are registered callbacks invoked upon Solace
Connectivity events such as disconnect/reconnecting/reconnected
events. See also type ConEvent.
*/


/*
The Callbacks struct contains callback function references for all
Solace events and is passed into the Init function.

- MsgCB: Event callback for message events; this is invoked for all message 
transports, Guaranteed and Direct.
     
- ErrCB: Event callback for error events.
     
- PubCB: Event callback for publisher events. Typically this is used
to handle asynchronous message-acks from the broker when the
publisher is configured to stream Guaranteed messages (rather than
publish and wait for the ack synchronously like a JMS persistent
publisher). Possible events include message ACK and REJECT events.
     
- ConCB: Event callback for connectivity events. This allows your
applications to be notified when the connection is lost or restored via
the Solace API's auto-reconnect logic.

See Init documentation and type documentation for callback function
types MsgHandler, ErrHandler, PubHandler, and ConHandler, as well as
their appropriate event type MsgEvent, ErrEvent, ConEvent, PubEvent.
*/


/*
Init initializes the underlying Solace API and returns a SESSION instance
to be passed into all further invocations of the gosol API functions to
identify the proper Solace SESSION for each call. This allows applications
to open multiple independent Solace sessions.

To Initialize a gosol SESSION, pass in references to all application
event-callback functions, which are always invoked asynchronously on a
background Solace thread allocated in the native Solace library. You do
not need to implement them all, see the documentation for type Callbacks
*/


/*
Connect establishes a TCP connection to a Solace message broker for the provided
SESSION with the session properties provided in the propsfile. For details about
Solace session properties, consult the Solace API reference documentation.
*/


/*
Connect establishes a TCP connection to a Solace message broker for the provided
SESSION with the session properties provided in the propsfile. For details about
Solace session properties, consult the Solace API reference documentation.
*/

/*
Disconnect gracefully shuts down a TCP connection to a Solace message broker 
for the provided SESSION.
*/


/*
SendDirect publishes a message to a Solace message broker via
an established SESSION connection using DIRECT (non-guaranteed)
transport. The message is published to the provided topic; the message
payload is expected to be Serialized by the calling application and
provided to gosol as a binary data buffer.
*/


/*
SendPersistent publishes a message to a Solace message broker via ann
established SESSION connection usinng Guaranteed transport. The message
can be published either to a topic or to a queue, distinguished via the
desttype parameter (TOPIC | QUEUE).
*/
/*
SendPersistentStreaming publishes a message to a Solace message broker
via ann established SESSION connection usinng Guaranteed transport (see
SendPersistent documentation). Additionally, SendPersistentStreaming
supports the transmission of additional Correlation data that will be
returned toe the application via the PubHandler callback function for
purposes of correlating streaming message acknowledgements from the
Solace message broker with the original sent messages.
*/


/*
SubscribeTopic subscribes to messages published to a Solace message
broker via an established SESSION connection with a topic string
representing the consumer's filtering interests. Messages will be consumed
via DIRECT (non-guaranteed) transport.
*/


/*
UnsubscribeTopic unsubscribes from a previously subscribed topic to a Solace 
message broker via an established SESSION connection.
*/


/*
BindQueue establishes a connection to a named queue on a Solace
message broker via an established SESSION connection. All matching
cache messages will be delivered to the application via the MsgHandler
callback function. The binding allows the consumer to consume messages
from the queue with 2 optional forward modes (STORE_FWD, CUT_THRU) and
2 optional ack modes (AUTO_ACK, MANUAL_ACK). In the case of MANUAL_ACK,
the application is required to acknowledge messages back to the Solace
message router for removal from the queue via the gosol.AckMsg function.
*/


/*
UnbindQueue releases a connection to a named queue on a Solace message
broker via an establed SESSION connection. The released binding allows
other consumers to continue consuming messages from the queue, and the
queue will not be destroyed.
*/


/*
AckMsg allows applications to Acknowledge messages that were delivered
over Guaranteed transport from a via an established SESSION connection
to a Solace message broker. This is only required for messages that were
delivered as a result of binding to a Solace message queue via BindQueue
with MANUAL_ACK ack-mode.
*/


/*
CacheReq requests messages from a SolCache instance identified via
cach_name and matching a filtering subscription interest expressed via
the topic_sub parameter. All matching cache messages will be delivered
to the application via the MsgHandler callback function.
*/


/*
Cut-Through Messaging allows for the delivery of Guaranteed messages with very low latency from Solace PubSub+ to consumers. 
This is done by using the low-latency, Direct Messaging data path for the bulk of the message flow, 
while also relying on the Guaranteed Messaging data path for message recovery in the event of a message loss. 
Cut-Through Messaging is not supported when the corresponding queue or topic endpoint is configured to respect message priority values
*/


	// - MsgCB: Event callback for message events; this is invoked for all message 
	// transports, Guaranteed and Direct.
	     
	// - ErrCB: Event callback for error events.
	     
	// - ConCB: Event callback for connectivity events. This allows your
	// applications to be notified when the connection is lost or restored via
	// the Solace API's auto-reconnect logic.

	// - PubCB: Event callback for publisher events. Typically this is used
	// to handle asynchronous message-acks from the broker when the
	// publisher is configured to stream Guaranteed messages (rather than
	// publish and wait for the ack synchronously like a JMS persistent
	// publisher). Possible events include message ACK and REJECT events.
