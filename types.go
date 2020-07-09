package solace

import "unsafe"

type FWD_MODE uint32
type ACK_MODE uint32
type EVENT_TYPE uint32
type CONNECTIVITY_STATUS uint32
type DESTINATION_TYPE uint32

const (
    // Fwd-Mode enumeration for guaranteed message delivery to consumers; 
    // see argument to func BindQueue()
    FWD_MODE_STORE_FWD FWD_MODE   = 1
    FWD_MODE_CUT_THRU  FWD_MODE   = 2
)

const (
    // Ack-Mode enumeration for consumer acknowledgement of guaranteed
    // message; in the case of MANUAL_ACK, the consuming application
    // is responsible for calling the AckMsg() function for messages
    // it has succesfully consumed/processed and thus does not require
    // redelivery. See argument to func BindQueue().
    ACK_MODE_AUTO     ACK_MODE = 1
    ACK_MODE_MANUAL   ACK_MODE = 2
)

const (
    // Publisher Event Type enumeration for publisher-related guaranteed
    // message events; possible events are message ACK and REJECT. See
    // PubEvent struct member Type.
    EVENT_TYPE_ACK          EVENT_TYPE = 1
    EVENT_TYPE_REJECT       EVENT_TYPE = 2
)

const (
    // Connectivity Event Type enumeration for session-related
    // connectivity events; possible events triggered include connection
    // UP/DOWN, RECONNECTING and RECONNECTED. See ConEvent struct member
    // Type.
    CONNECTIVITY_STATUS_UP           CONNECTIVITY_STATUS = 1
    CONNECTIVITY_STATUS_RECONNECTING CONNECTIVITY_STATUS = 2
    CONNECTIVITY_STATUS_RECONNECTED  CONNECTIVITY_STATUS = 3
    CONNECTIVITY_STATUS_DOWN         CONNECTIVITY_STATUS = 4
)

const (
    // Destination enumeration for published messages over any transport;
    // possible destinations are TOPIC or QUEUE. see SendDirect documentation
    // and MsgEvent member DestType.
    DESTINATION_TYPE_TOPIC        DESTINATION_TYPE = 1
    DESTINATION_TYPE_QUEUE        DESTINATION_TYPE = 2
    DESTINATION_TYPE_NONE         DESTINATION_TYPE = 3
)


/*
ErrEvent is passed to a registered ErrHandler callback for each Solace 
error-event received
*/
type ErrorEvent  struct {
	Session       uint64
	FunctionName  string
	ReturnCode    int
	RCStr         string
	SubCode       int
	SCStr         string
	ResponseCode  int
	ErrorStr      string
}


/*
PubEvent is passed to a registered PubHandler callback for each
Solace publisher-event received such as message ACK or REJECT events; for
streaming publishers it can include publisher-provided correlation data.
*/
type PublisherEvent struct {
	Session         uint64
	Type            string
	CorrelationData unsafe.Pointer
}


/*
ConEvent is passed to a registered ConHandler callback for each Solace
connectivity-event received such as disconnect/reconnecting/reconnected
events
*/
type ConnectionEvent struct {
	Session   uint64
	Type      string
}

/*
MessageEvent is passed to a registered MsgHandler callback for each Solace
message-event received, for either Guaranteed or Direct transports
*/
type MessageEvent struct {
	Session            uint64

	DestinationType    string
	Destination        string

	ReplyToDestinationType    string
	ReplyToDestination        string

	Flow               uint64
	MessageID          uint64

	MessageType        string
	BinaryPayload      []byte
	BinaryPayloadLen   uint

	CorrelationID      string

	RequestID          int
	Redelivered        bool
	Discard            bool
	UserProperties     string
}

type ISolace interface {
	
	GetMessageChannel() chan MessageEvent
	GetErrorChannel() chan ErrorEvent
	GetConnectionChannel() chan ConnectionEvent
	GetPublisherChannel() chan PublisherEvent
	
	NotifyOnMessage(func(msg MessageEvent))
	NotifyOnConnection(callback func(msg ConnectionEvent))
	NotifyOnError(callback func(msg ErrorEvent))
	NotifyOnPublisher(callback func(msg PublisherEvent))

	Connect(host string, vpn string, user string, pass string, windowsize string)
	
	Disconnect()
	
	Subscribe(topic string)
	Unsubscribe(topic string)
	
	Publish(
		destinationType DESTINATION_TYPE, target string, 
		replyToDestType DESTINATION_TYPE, replyTo string, 
		messageType string, 
		payload []byte, 
		userProperties map[string]interface{}, 
		correlationID   string,
		correlationData int,
	)
	
	SubscribeQueue(queueName string, mode ACK_MODE)
	UnsubscribeQueue(queueName string)
	
	Ack(flow uint64, msgID uint64)
}