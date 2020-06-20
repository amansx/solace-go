package main

import "C"
import "unsafe"
import "gosol"
import "gosol_t"

var singleton *Solace

type Solace struct {
	sess gosol.SESSION
}

func (this *Solace) Init(onMessage gosol_t.MessageHandler, onError gosol_t.ErrorHandler, onConnection gosol_t.ConnectionEventHandler, onPublisherEvent gosol_t.PublisherEventHandler) {
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

	this.sess = gosol.Init(
		gosol.MsgHandler(onMessage),
		gosol.ErrHandler(onError),
		gosol.ConHandler(onConnection), 
		gosol.PubHandler(onPublisherEvent),
	)
}

func (this Solace) Connect(host string, vpn string, user string, pass string, windowsize string) {
	gosol.ConnectWithParams(this.sess, host, vpn, user, pass, windowsize)
}

func (this Solace) Disconnect() {
	gosol.Disconnect(this.sess)
}

func (this Solace) Subscribe(topic string) {
	gosol.SubscribeTopic(this.sess, topic)
}

func (this Solace) Unsubscribe(topic string) {
	gosol.UnsubscribeTopic(this.sess, topic)
}

func (this Solace) Publish(topicName string, payload []byte) {
	payloadptr := unsafe.Pointer(&payload[0])
	payloadLen := len(payload)
	gosol.SendPersistent(this.sess, topicName, gosol.TOPIC, payloadptr, payloadLen)
}

func (this Solace) PublishDirect(topic string, payload []byte) {
	payloadptr := unsafe.Pointer(&payload[0])
	payloadLen := len(payload)
	gosol.SendDirect(this.sess, topic, payloadptr, payloadLen)
}

func (this Solace) PublishToQueue(queueName string, payload []byte) {
	payloadptr := unsafe.Pointer(&payload[0])
	payloadLen := len(payload)
	gosol.SendPersistent(this.sess, queueName, gosol.QUEUE, payloadptr, payloadLen)
}

/*
Cut-Through Messaging allows for the delivery of Guaranteed messages with very low latency from Solace PubSub+ to consumers. 
This is done by using the low-latency, Direct Messaging data path for the bulk of the message flow, 
while also relying on the Guaranteed Messaging data path for message recovery in the event of a message loss. 
Cut-Through Messaging is not supported when the corresponding queue or topic endpoint is configured to respect message priority values
*/
func (this Solace) BindQueue(queueName string) {
	gosol.BindQueue(this.sess, queueName, gosol.STORE_FWD, gosol.AUTO_ACK)
}

func (this Solace) BindQueueManualAck(queueName string) {
	gosol.BindQueue(this.sess, queueName, gosol.STORE_FWD, gosol.MANUAL_ACK)
}

func (this Solace) UnbindQueue(queueName string) {
	gosol.UnbindQueue(this.sess, queueName)
}

func (this Solace) Ack(flow uint64, msgID uint64) {
	gosol.AckMsg(this.sess, flow, msgID)
}


func GoSolace(onMessage gosol_t.MessageHandler, onError gosol_t.ErrorHandler, onConnection gosol_t.ConnectionEventHandler, onPublisherEvent gosol_t.PublisherEventHandler) gosol_t.Solace {
	if (singleton == nil) {
		singleton = &Solace{}
		singleton.Init(onMessage, onError, onConnection, onPublisherEvent)
	}
	return singleton
}