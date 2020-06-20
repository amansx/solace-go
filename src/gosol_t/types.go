package gosol_t

import "unsafe"

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
	Flow               uint64
	MessageID          uint64
	Buffer             []byte
	BufferLen          uint
	RequestID          int
	Redelivered        bool
	Discard            bool
	UserProperties     string
}

type MessageHandler func(MessageEvent)
type ErrorHandler func(ErrorEvent)
type PublisherEventHandler func(PublisherEvent)
type ConnectionEventHandler func(ConnectionEvent)

type Solace interface {
	Init(MessageHandler, ErrorHandler, ConnectionEventHandler, PublisherEventHandler)
	Connect(host string, vpn string, user string, pass string, windowsize string)
	Disconnect()
	Subscribe(topic string)
	Unsubscribe(topic string)
	Publish(topicName string, payload []byte)
	PublishDirect(topicName string, payload []byte)
	PublishToQueue(queueName string, payload []byte)
	BindQueue(queueName string)
	BindQueueManualAck(queueName string)
	UnbindQueue(queueName string)
	Ack(flow uint64, msgID uint64)
}