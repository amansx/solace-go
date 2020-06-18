package main

import "gosol"

type Solace struct {
	sess gosol.SESSION
}

func (this *Solace) Init() {
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

	this.sess := gosol.Init(
		gosol.MsgHandler(on_msg),
		gosol.ErrHandler(on_err), 
		nil, 
		nil
	)	
}

func (this Solace) Connect() {
	gosol.Connect(this.sess, propertiesFilePath)
}

func (this Solace) Disconnect() {
	gosol.Disconnect(this.sess)
}

func (this Solace) Subscribe(topic String) {
	gosol.SubscribeTopic(this.sess, topic)
}

func (this Solace) Unsubscribe(topic String) {
	gosol.UnsubscribeTopic(this.sess, topic)
}

func (this Solace) PublishDirect(topic String, payload []byte) {
	payloadptr := unsafe.Pointer(&bytes[0])
	payloadLen := len(payload)
	gosol.SendDirect(this.sess, topic, payloadptr, payloadLen)
}