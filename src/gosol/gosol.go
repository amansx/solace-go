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
package gosol

// #include "sol_api.h"
// #include <stdlib.h>
//
// extern void gosol_on_msg(SOLHANDLE h, message_event* msg);
// extern void gosol_on_err(SOLHANDLE h, error_event* err);
// extern void gosol_on_pub(SOLHANDLE h, publisher_event* pub);
// extern void gosol_on_con(SOLHANDLE h, connectivity_event* con);
// 
// static void on_msg(SOLHANDLE h, message_event* msg) {
//	gosol_on_msg( h, msg );
// }
// static void on_err(SOLHANDLE h, error_event* err) {
//	gosol_on_err( h, err );
// }
// static void on_pub(SOLHANDLE h, publisher_event* pub) {
//	gosol_on_pub( h, pub );
// }
// static void on_con(SOLHANDLE h, connectivity_event* con) {
//	gosol_on_con( h, con );
// }
// static SOLHANDLE gosol_init(void* gohandlers) {
//	callback_functions cbs;
//	cbs.err_cb = on_err;
//	cbs.msg_cb = on_msg;
//	cbs.pub_cb = on_pub;
//	cbs.con_cb = on_con;
//	cbs.user_data = gohandlers;
// 	return sol_init( on_msg, on_err, on_pub, on_con, gohandlers );
// }
import "C"
import "unsafe"

// Session handle is return from Init(), and passed into all other functions
type SESSION uint64

const (
    // Fwd-Mode enumeration for guaranteed message delivery to consumers; 
    // see argument to func BindQueue()
    STORE_FWD    = 1
    CUT_THRU     = 2
    
    // Ack-Mode enumeration for consumer acknowledgement of guaranteed
    // message; in the case of MANUAL_ACK, the consuming application
    // is responsible for calling the AckMsg() function for messages
    // it has succesfully consumed/processed and thus does not require
    // redelivery. See argument to func BindQueue().
    AUTO_ACK     = 1
    MANUAL_ACK   = 2
    
    // Publisher Event Type enumeration for publisher-related guaranteed
    // message events; possible events are message ACK and REJECT. See
    // PubEvent struct member Type.
    ACK          = 1
    REJECT       = 2
    
    // Connectivity Event Type enumeration for session-related
    // connectivity events; possible events triggered include connection
    // UP/DOWN, RECONNECTING and RECONNECTED. See ConEvent struct member
    // Type.
    UP           = 1
    RECONNECTING = 2
    RECONNECTED  = 3
    DOWN         = 4
    
    // Destination enumeration for published messages over any transport;
    // possible destinations are TOPIC or QUEUE. see SendDirect documentation
    // and MsgEvent member DestType.
    TOPIC        = 1
    QUEUE        = 2
    NONE         = 3
)


/*
MsgEvent is passed to a registered MsgHandler callback for each Solace
message-event received, for either Guaranteed or Direct transports
*/
type MsgEvent  struct {
	DestType    int
	Destination string
	Flow        uint64
	MsgId       uint64
	Buffer      unsafe.Pointer
	BufLen      uint
	ReqId       int
	Redelivered bool
	Discard     bool
}


/*
ErrEvent is passed to a registered ErrHandler callback for each Solace 
error-event received
*/
type ErrEvent  struct {
	FNName    string
	RetCode   int
	RCStr     string
	SubCode   int
	SCStr     string
	RespCode  int
	ErrStr    string
}


/*
PubEvent is passed to a registered PubHandler callback for each
Solace publisher-event received such as message ACK or REJECT events; for
streaming publishers it can include publisher-provided correlation data.
*/
type PubEvent struct {
	Type            int
	CorrelationData unsafe.Pointer
}


/*
ConEvent is passed to a registered ConHandler callback for each Solace
connectivity-event received such as disconnect/reconnecting/reconnected
events
*/
type ConEvent struct {
	Type      int
}


/*
MsgHandler functions are registered callbacks invoked upon receipt
of Solace messages for all transports (Guaranteed and Direct). See also
type MsgEvent.
*/
type MsgHandler func(sess SESSION, msg *MsgEvent)

/*
ErrHandler functions are registered callbacks invoked upon Solace
errors. See also type ErrEvent.
*/
type ErrHandler func(sess SESSION, err *ErrEvent)

/*
PubHandler functions are registered callbacks invoked upon Solace
Publisher events such as message ACK or REJECT events; for streaming
publishers it can include publisher-provided correlation data.. See also
type PubEvent.
*/
type PubHandler func(sess SESSION, pub *PubEvent)

/*
ConHandler functions are registered callbacks invoked upon Solace
Connectivity events such as disconnect/reconnecting/reconnected
events. See also type ConEvent.
*/
type ConHandler func(sess SESSION, con *ConEvent)

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
type Callbacks struct {
	MsgCB MsgHandler
	ErrCB ErrHandler
	ConCB ConHandler
	PubCB PubHandler
}


func i2b(i C.int) bool {
	if (i > 0) {
		return true
	}
	return false
}

//export gosol_on_msg
func gosol_on_msg(h SESSION, m *C.struct_message_event) {
	ctx := (*Callbacks)(m.user_data)

	if ctx.MsgCB != nil {
		evt            := new(MsgEvent)
		evt.DestType    = int(m.desttype)
		evt.Destination = C.GoString(m.destination)
		evt.Flow        = uint64(m.flow)
		evt.MsgId       = uint64(m.id)
		evt.Buffer      = unsafe.Pointer(m.buffer)
		evt.BufLen      = uint(m.buflen)
		evt.ReqId       = int(m.req_id)
		evt.Redelivered = i2b(m.redelivered_flag)
		evt.Discard     = i2b(m.discard_flag)

		ctx.MsgCB( SESSION(h), evt )
	}
}

//export gosol_on_err
func gosol_on_err(h SESSION, e *C.struct_error_event) {
	ctx := (*Callbacks)(e.user_data)

	if ctx.ErrCB != nil {
		evt         := new(ErrEvent)
		evt.FNName   = C.GoString(e.fn_name)
		evt.RetCode  = int(e.return_code)
		evt.RCStr    = C.GoString(e.rc_str)
		evt.SubCode  = int(e.sub_code)
		evt.SCStr    = C.GoString(e.sc_str)
		evt.RespCode = int(e.resp_code)
		evt.ErrStr   = C.GoString(e.err_str)

		ctx.ErrCB( SESSION(h), evt )
	}
}

//export gosol_on_pub
func gosol_on_pub(h SESSION, p *C.struct_publisher_event) {
	ctx := (*Callbacks)(p.user_data)

	if ctx.PubCB != nil {
		evt                := new(PubEvent)
		evt.Type            = int(p._type)
		evt.CorrelationData = unsafe.Pointer(p.correlation_data)

		ctx.PubCB( SESSION(h), evt )
	}
}

//export gosol_on_con
func gosol_on_con(h SESSION, c *C.struct_connectivity_event) {
	ctx := (*Callbacks)(c.user_data)

	if ctx.ConCB != nil {
		evt      := new(ConEvent)
		evt.Type  = int(c._type)

		ctx.ConCB( SESSION(h), evt )
	}
}


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
func Init(MsgCB MsgHandler, ErrCB ErrHandler, ConCB ConHandler, PubCB PubHandler) SESSION {
	cbs := Callbacks { MsgCB, ErrCB, ConCB, PubCB }
	sess := C.gosol_init( (unsafe.Pointer)(unsafe.Pointer(&cbs)) )
	return SESSION(sess)
}

/*
Connect establishes a TCP connection to a Solace message broker for the provided
SESSION with the session properties provided in the propsfile. For details about
Solace session properties, consult the Solace API reference documentation.
*/
func Connect(sess SESSION, propsfile string) int {
	// C.CString allocates new memory on the heap, so need to cleanup afterwards
	cprops := C.CString(propsfile)
	rc     := int( C.sol_connect( C.SOLHANDLE(sess), cprops ) )
	C.free( unsafe.Pointer(cprops) )
	return rc
}

/*
Disconnect gracefully shuts down a TCP connection to a Solace message broker 
for the provided SESSION.
*/
func Disconnect(sess SESSION) int {
	return int( C.sol_disconnect( C.SOLHANDLE(sess) ) )
}


/*
SendDirect publishes a message to a Solace message broker via
an established SESSION connection using DIRECT (non-guaranteed)
transport. The message is published to the provided topic; the message
payload is expected to be Serialized by the calling application and
provided to gosol as a binary data buffer.
*/
func SendDirect(sess SESSION, topic string, buffer unsafe.Pointer, buflen int) int {
	ctopic := C.CString(topic)
	rc     := int( C.sol_send_direct( C.SOLHANDLE(sess), ctopic, buffer, C.int(buflen) ) )
	C.free( unsafe.Pointer(ctopic) )
	return rc
}

/*
SendPersistent publishes a message to a Solace message broker via ann
established SESSION connection usinng Guaranteed transport. The message
can be published either to a topic or to a queue, distinguished via the
desttype parameter (TOPIC | QUEUE).
*/
func SendPersistent(sess SESSION, dest string, desttype uint32, buffer unsafe.Pointer, buflen int) int {
	cdest := C.CString(dest)
	rc    := int( C.sol_send_persistent( C.SOLHANDLE(sess), cdest, desttype, buffer, C.int(buflen), nil, 0 ) )
	C.free( unsafe.Pointer(cdest) )
	return rc
}


/*
SendPersistentStreaming publishes a message to a Solace message broker
via ann established SESSION connection usinng Guaranteed transport (see
SendPersistent documentation). Additionally, SendPersistentStreaming
supports the transmission of additional Correlation data that will be
returned toe the application via the PubHandler callback function for
purposes of correlating streaming message acknowledgements from the
Solace message broker with the original sent messages.
*/
func SendPersistentStreaming(sess SESSION, dest string, desttype uint32, buffer unsafe.Pointer, buflen int, corrptr unsafe.Pointer, corrlen int) int {
	cdest := C.CString(dest)
	rc    := int( C.sol_send_persistent( C.SOLHANDLE(sess), cdest, desttype, buffer, C.int(buflen), corrptr, C.int(corrlen) ) )
	C.free( unsafe.Pointer(cdest) )
	return rc
}


/*
SubscribeTopic subscribes to messages published to a Solace message
broker via an established SESSION connection with a topic string
representing the consumer's filtering interests. Messages will be consumed
via DIRECT (non-guaranteed) transport.
*/
func SubscribeTopic(sess SESSION, topic string) int {
	ctopic := C.CString(topic)
	rc     := int( C.sol_subscribe_topic( C.SOLHANDLE(sess), ctopic ) )
	C.free( unsafe.Pointer(ctopic) )
	return rc
}


/*
UnsubscribeTopic unsubscribes from a previously subscribed topic to a Solace 
message broker via an established SESSION connection.
*/
func UnsubscribeTopic(sess SESSION, topic string) int {
	ctopic := C.CString(topic)
	rc     := int( C.sol_unsubscribe_topic( C.SOLHANDLE(sess), ctopic ) )
	C.free( unsafe.Pointer(ctopic) )
	return rc
}

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
func BindQueue(sess SESSION, queue string, fwdmode uint32, ackmode uint32) int {
	cqueue := C.CString(queue)
	rc     := int( C.sol_bind_queue( C.SOLHANDLE(sess), cqueue, fwdmode, ackmode ) )
	C.free( unsafe.Pointer(cqueue) )
	return rc
}


/*
UnbindQueue releases a connection to a named queue on a Solace message
broker via an establed SESSION connection. The released binding allows
other consumers to continue consuming messages from the queue, and the
queue will not be destroyed.
*/
func UnbindQueue(sess SESSION, queue string) int {
	cqueue := C.CString(queue)
	rc     := int( C.sol_unbind_queue( C.SOLHANDLE(sess), cqueue ) )
	C.free( unsafe.Pointer(cqueue) )
	return rc
}

/*
AckMsg allows applications to Acknowledge messages that were delivered
over Guaranteed transport from a via an established SESSION connection
to a Solace message broker. This is only required for messages that were
delivered as a result of binding to a Solace message queue via BindQueue
with MANUAL_ACK ack-mode.
*/
func AckMsg(sess SESSION, flow uint64, msgid uint64) int {
	return int( C.sol_ack_msg( C.SOLHANDLE(sess), C.FLOWHANDLE(flow), C.SOLMSGID(msgid) ) )
}

/*
CacheReq requests messages from a SolCache instance identified via
cach_name and matching a filtering subscription interest expressed via
the topic_sub parameter. All matching cache messages will be delivered
to the application via the MsgHandler callback function.
*/
func CacheReq(sess SESSION, cache_name string, topic_sub string , request_id int) int {
	ccache_name := C.CString(cache_name)
	ctopic_sub  := C.CString(topic_sub)
	rc          := int( C.sol_cache_req( C.SOLHANDLE(sess), ccache_name, ctopic_sub, C.int(request_id) ) )
	C.free( unsafe.Pointer(ccache_name) )
	C.free( unsafe.Pointer(ctopic_sub) )
	return rc
}


