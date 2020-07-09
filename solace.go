package main

import "C"
import "unsafe"
import "encoding/json"

type Solace struct {
	sess SESSION
	messageEventChannel    chan MessageEvent
	errorEventChannel      chan ErrorEvent
	connectionEventChannel chan ConnectionEvent
	publisherEventChannel  chan PublisherEvent
}

func (this *Solace) GetMessageChannel() chan MessageEvent {
	return this.messageEventChannel
}
func (this *Solace) GetErrorChannel() chan ErrorEvent {
	return this.errorEventChannel
}
func (this *Solace) GetConnectionChannel() chan ConnectionEvent {
	return this.connectionEventChannel
}
func (this *Solace) GetPublisherChannel() chan PublisherEvent {
	return this.publisherEventChannel
}

func (this Solace) Connect(host string, vpn string, user string, pass string, windowsize string) {
	ConnectWithParams(this.sess, host, vpn, user, pass, windowsize)
}

func (this Solace) NotifyOnMessage(callback func(msg MessageEvent)) {
	go func() {
		for {
			select {
			case msg := <-this.messageEventChannel:
			callback(msg)
			}
		}
	}()
}
func (this Solace) NotifyOnConnection(callback func(msg ConnectionEvent)) {
	go func() {
		for {
			select {
			case con := <-this.connectionEventChannel:
			callback(con)
			}
		}
	}()
}
func (this Solace) NotifyOnError(callback func(msg ErrorEvent)) {
	go func() {
		for {
			select {
			case err := <-this.errorEventChannel:
			callback(err)
			}
		}
	}()
}
func (this Solace) NotifyOnPublisher(callback func(msg PublisherEvent)) {
	go func() {
		for {
			select {
			case pub := <-this.publisherEventChannel:
			callback(pub)
			}
		}
	}()
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

func (this Solace) Publish(
	destinationType DESTINATION_TYPE, target string, 
	replyToDestType DESTINATION_TYPE, replyTo string, 
	messageType string, 
	payload []byte, 
	userProperties map[string]interface{}, 
	correlationID string,
	correlationData int,
) {
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
	
	var payloadptr unsafe.Pointer
	payloadLen := len(payload)
	if payloadLen > 0 {
		payloadptr = unsafe.Pointer(&payload[0])
	} else {
		payloadptr = nil
	}
	

	if (destinationType == DESTINATION_TYPE_NONE) {
		SendDirect(this.sess, target, payloadptr, payloadLen, string(propsJson))

	} else {

		corrdPtr := C.int(correlationData)

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

			correlationID,

			unsafe.Pointer(&corrdPtr), 
			int(unsafe.Sizeof(corrdPtr)),
		)
	}
	
}

func (this Solace) SubscribeQueue(queueName string, mode ACK_MODE) {
	BindQueue(this.sess, queueName, FWD_MODE_STORE_FWD, mode)
}

func (this Solace) UnsubscribeQueue(queueName string) {
	UnbindQueue(this.sess, queueName)
}

func (this Solace) Ack(flow uint64, msgID uint64) {
	AckMsg(this.sess, flow, msgID)
}


const BackPressureLimit = 5000
func NewSolace() ISolace {
	solace := &Solace{
		messageEventChannel: make(chan MessageEvent, BackPressureLimit),
		errorEventChannel: make(chan ErrorEvent, BackPressureLimit),
		connectionEventChannel: make(chan ConnectionEvent, BackPressureLimit),
		publisherEventChannel: make(chan PublisherEvent, BackPressureLimit),
	}
	solace.sess = Init(solace)
	return solace
}