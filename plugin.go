// +build core

package main

import "C"
import "unsafe"
import "encoding/json"
import "internal/solace"

var singleton *Solace

type Solace struct {
	sess SESSION
}

func (this *Solace) Init(onMessage solace.MessageHandler, onError solace.ErrorHandler, onConnection solace.ConnectionEventHandler, onPublisherEvent solace.PublisherEventHandler) {
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

	this.sess = Init(
		MsgHandler(onMessage),
		ErrHandler(onError),
		ConHandler(onConnection), 
		PubHandler(onPublisherEvent),
	)
}

func (this Solace) Connect(host string, vpn string, user string, pass string, windowsize string) {
	ConnectWithParams(this.sess, host, vpn, user, pass, windowsize)
}

func (this Solace) Disconnect() {
	Disconnect(this.sess)
}

func (this Solace) Subscribe(topic string) {
	SubscribeTopic(this.sess, topic)
}

func (this Solace) Unsubscribe(topic string) {
	UnsubscribeTopic(this.sess, topic)
}

func (this Solace) Publish(destinationType solace.DESTINATION_TYPE, target string, replyToDestType solace.DESTINATION_TYPE, replyTo string, messageType string, payload []byte, userProperties map[string]interface{}, corelationID int) {
	props := map[string]map[string]interface{}{}
	for k, v := range userProperties {
		switch v.(type){
			case int8: 
				if _, ok := props["int8"]; !ok {
					props["int8"] = map[string]interface{}{}
				}
				props["int8"][k] = v.(interface{})
			case int16: 
				if _, ok := props["int16"]; !ok {
					props["int16"] = map[string]interface{}{}
				}
				props["int16"][k] = v.(interface{})
			case int32: 
				if _, ok := props["int32"]; !ok {
					props["int32"] = map[string]interface{}{}
				}
				props["int32"][k] = v.(interface{})
			case int, int64: 
				if _, ok := props["int64"]; !ok {
					props["int64"] = map[string]interface{}{}
				}
				props["int64"][k] = v.(interface{})
			case uint8: 
				if _, ok := props["uint8"]; !ok {
					props["uint8"] = map[string]interface{}{}
				}
				props["uint8"][k] = v.(interface{})
			case uint16: 
				if _, ok := props["uint16"]; !ok {
					props["uint16"] = map[string]interface{}{}
				}
				props["uint16"][k] = v.(interface{})
			case uint32: 
				if _, ok := props["uint32"]; !ok {
					props["uint32"] = map[string]interface{}{}
				}
				props["uint32"][k] = v.(interface{})
			case uint64:
				if _, ok := props["uint64"]; !ok {
					props["uint64"] = map[string]interface{}{}
				}
				props["uint64"][k] = v.(interface{})
			case bool:
				if _, ok := props["bool"]; !ok {
					props["bool"] = map[string]interface{}{}
				}
				props["bool"][k] = v.(interface{})
			case string:
				if _, ok := props["string"]; !ok {
					props["string"] = map[string]interface{}{}
				}
				props["string"][k] = v.(interface{})
		}
	}

	propsJson, _ := json.Marshal(props)
	
	payloadptr := unsafe.Pointer(&payload[0])
	payloadLen := len(payload)
	

	if (destinationType == solace.DESTINATION_TYPE_NONE) {
		SendDirect(this.sess, target, payloadptr, payloadLen, string(propsJson))

	} else {

		coordID := C.int(corelationID)

		SendPersistentWithCorelation(
			this.sess, 

			target, 
			destinationType, 
			
			replyTo,
			replyToDestType,

			messageType,

			payloadptr, 
			payloadLen, 

			string(propsJson),

			unsafe.Pointer(&coordID), 
			int(unsafe.Sizeof(coordID)),
		)
	}
	
}

/*
Cut-Through Messaging allows for the delivery of Guaranteed messages with very low latency from Solace PubSub+ to consumers. 
This is done by using the low-latency, Direct Messaging data path for the bulk of the message flow, 
while also relying on the Guaranteed Messaging data path for message recovery in the event of a message loss. 
Cut-Through Messaging is not supported when the corresponding queue or topic endpoint is configured to respect message priority values
*/
func (this Solace) SubscribeQueue(queueName string, mode solace.ACK_MODE) {
	BindQueue(this.sess, queueName, solace.FWD_MODE_STORE_FWD, mode)
}

func (this Solace) UnsubscribeQueue(queueName string) {
	UnbindQueue(this.sess, queueName)
}

func (this Solace) Ack(flow uint64, msgID uint64) {
	AckMsg(this.sess, flow, msgID)
}


var InitSolaceWithCallback solace.InitSolaceWithCallback = func (onMessage solace.MessageHandler, onError solace.ErrorHandler, onConnection solace.ConnectionEventHandler, onPublisherEvent solace.PublisherEventHandler) solace.ISolace {
	if (singleton == nil) {
		singleton = &Solace{}
		singleton.Init(onMessage, onError, onConnection, onPublisherEvent)
	}
	return singleton
}