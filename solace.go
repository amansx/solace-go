package solace

//#cgo amd64 386 CFLAGS: -DPROVIDE_LOG_UTILITIES -DSOLCLIENT_CONST_PROPERTIES
//#cgo CFLAGS: -I. -I${SRCDIR}/includes/solwrap/include -fPIC -m64 -g
//#cgo linux CFLAGS: -D_REENTRANT -D_LINUX_X86_64
//#cgo linux LDFLAGS: -L${SRCDIR}/lib.linux -lsolwrap -lsolclient -lpthread -lrt -lstdc++
//#cgo darwin CFLAGS: -D_REENTRANT
//#cgo darwin LDFLAGS: -L${SRCDIR}/lib.osx -lsolwrap -lsolclient -lcrypto -lssl -lgssapi_krb5 -lpthread -lstdc++
//#cgo windows CFLAGS: -DSOL_WRAP_API -DPROVIDE_LOG_UTILITIES -DSOLCLIENT_CONST_PROPERTIES -DWIN32 -DNDEBUG -D_WINDOWS -D_USRDLL
//#cgo windows LDFLAGS: -L${SRCDIR}/lib.win -llibsolwrap -llibsolclient -static-libstdc++ -static-libgcc -lstdc++
//#include <stdlib.h>
//#include "sol_api.h"
//#include "solace.c"
import "C"
import "unsafe"
import "encoding/json"
import cptr "github.com/mattn/go-pointer"

type SESSION uint64

//export gosol_on_msg
func gosol_on_msg(h SESSION, m *C.struct_message_event) {
	ctx := cptr.Restore(m.user_data).(*Solace)
	if n := ctx.GetMessageChannel(); n != nil {

		evt                        := MessageEvent{}
		evt.Session                 = uint64(h)
		evt.DestinationType         = destinationTypeToString(int(m.desttype))
		evt.Destination             = C.GoString(m.destination)
		evt.ReplyToDestinationType  = destinationTypeToString(int(m.replyToDestType))
		evt.ReplyToDestination      = C.GoString(m.replyTo)
		evt.Flow                    = uint64(m.flow)
		evt.MessageID               = uint64(m.id)
		evt.RequestID               = int(m.req_id)
		evt.Redelivered             = i2b(m.redelivered_flag)
		evt.Discard                 = i2b(m.discard_flag)
		evt.UserProperties          = C.GoString(m.user_properties)
		evt.CorrelationID           = C.GoString(m.correlationid)
		
		if m.buffer != nil {
			evt.BinaryPayload    = C.GoBytes(unsafe.Pointer(m.buffer), C.int(m.buflen))
			evt.BinaryPayloadLen = uint(m.buflen)
		}

		if m.application_message_type != nil {
			evt.MessageType  = C.GoString(m.application_message_type)
		} else {
			evt.MessageType  = ""
		}

		n <- evt

	}
}

//export gosol_on_err
func gosol_on_err(h SESSION, e *C.struct_error_event) {
	ctx := cptr.Restore(e.user_data).(*Solace)
	if n := ctx.GetErrorChannel(); n != nil {
		evt              := ErrorEvent{}
		evt.Session      = uint64(h)
		evt.FunctionName = C.GoString(e.fn_name)
		evt.ReturnCode   = int(e.return_code)
		evt.RCStr        = C.GoString(e.rc_str)
		evt.SubCode      = int(e.sub_code)
		evt.SCStr        = C.GoString(e.sc_str)
		evt.ResponseCode = int(e.resp_code)
		evt.ErrorStr     = C.GoString(e.err_str)
		n <- evt
	}
}

//export gosol_on_pub
func gosol_on_pub(h SESSION, p *C.struct_publisher_event) {
	ctx := cptr.Restore(p.user_data).(*Solace)
	if n := ctx.GetPublisherChannel(); n != nil {
		evt                := PublisherEvent{}
		evt.Session         = uint64(h)
		evt.Type            = publishedEventTypeToString(int(p._type))
		// evt.CorrelationData = unsafe.Pointer(p.correlation_data)
		n <- evt
	}
}

//export gosol_on_con
func gosol_on_con(h SESSION, c *C.struct_connectivity_event) {
	ctx := cptr.Restore(c.user_data).(*Solace)
	if n := ctx.GetConnectionChannel(); ctx != nil && n != nil {
		evt := ConnectionEvent{
			Session: uint64(h),
			Type   : connectionTypeToString(int(c._type)),
		}
		n <- evt
	}
}

func i2b(i C.int) bool {
	if (i > 0) {
		return true
	}
	return false
}

func destinationTypeToString(i int) string {
	switch i {
		case 1: return "TOPIC"
		case 2: return "QUEUE"
	}
	return "NONE"
}

func connectionTypeToString(i int) string {
	switch i {
		case 1: return "UP"
		case 2: return "RECONNECTING"
		case 3: return "RECONNECTED"
	}
	return "DOWN"
}

func publishedEventTypeToString(i int) string {
	switch i {
	case 1:
		return "ACK"
	}
	return "REJECT"
}

var messagePayloadChannel chan[]byte



func Init(solace *Solace) SESSION {
	sess := C.gosol_init(cptr.Save(solace))
	return SESSION(sess)
}

func Connect(sess SESSION, propsfile string) int {
	// C.CString allocates new memory on the heap, so need to cleanup afterwards
	cprops := C.CString(propsfile)
	rc     := int( C.sol_connect( C.SOLHANDLE(sess), cprops ) )
	C.free( unsafe.Pointer(cprops) )
	return rc
}

func ConnectWithParams(sess SESSION, host string, vpn string, user string, pass string, appName string, appDesc string, windowsize string) int {
	// C.CString allocates new memory on the heap, so need to cleanup afterwards
	chost    := C.CString(host)
	cvpn     := C.CString(vpn)
	cuser    := C.CString(user)
	cpass    := C.CString(pass)
	cappName := C.CString(appName)
	cappDesc := C.CString(appDesc)
	cws      := C.CString(windowsize)
	rc      := int( C.sol_connect_with_params( C.SOLHANDLE(sess), chost, cvpn, cuser, cpass, cappName, cappDesc, cws ))

	C.free(unsafe.Pointer(chost))
	C.free(unsafe.Pointer(cvpn))
	C.free(unsafe.Pointer(cuser))
	C.free(unsafe.Pointer(cpass))
	C.free(unsafe.Pointer(cappName))
	C.free(unsafe.Pointer(cappDesc))
	C.free(unsafe.Pointer(cws))

	return rc
}

func Disconnect(sess SESSION) int {
	return int( C.sol_disconnect( C.SOLHANDLE(sess) ) )
}


func SendDirect(sess SESSION, topic string, buffer unsafe.Pointer, buflen int, userprops string) int {
	ctopic     := C.CString(topic)
	cuserProps := C.CString(userprops)
	
	rc     := int( C.sol_send_direct( C.SOLHANDLE(sess), ctopic, buffer, C.int(buflen), cuserProps ) )
	
	C.free( unsafe.Pointer(ctopic) )
	C.free( unsafe.Pointer(cuserProps) )
	return rc
}

func SendPersistentWithCorelation(
	sess SESSION, 

	dest string, 
	desttype DESTINATION_TYPE, 

	replyTo string, 
	replyToDestType DESTINATION_TYPE, 

	messageType string,

	buffer unsafe.Pointer, 
	buflen int, 

	userprops string,

	correlationid string,

	corrptr unsafe.Pointer, 
	corrlen int, 

) int {

	cdest          := C.CString(dest)
	creplyTo       := C.CString(replyTo)
	cuserProps     := C.CString(userprops)
	cmessageType   := C.CString(messageType)
	ccorrelationID := C.CString(correlationid)

	rc    := int( C.sol_send_persistent( 
		C.SOLHANDLE(sess), 
		
		cdest, 
		uint32(desttype), 

		creplyTo, 
		uint32(replyToDestType), 

		cmessageType, 

		buffer, 
		C.int(buflen), 

		cuserProps,

		ccorrelationID, 
		
		corrptr, 
		C.int(corrlen), 		
	
	))

	C.free( unsafe.Pointer(cdest) )
	C.free( unsafe.Pointer(creplyTo) )
	C.free( unsafe.Pointer(cuserProps) )
	C.free( unsafe.Pointer(cmessageType) )
	C.free( unsafe.Pointer(ccorrelationID) )

	return rc
}


func SubscribeTopic(sess SESSION, topic string) int {
	ctopic := C.CString(topic)
	rc     := int( C.sol_subscribe_topic( C.SOLHANDLE(sess), ctopic ) )
	C.free( unsafe.Pointer(ctopic) )
	return rc
}


func UnsubscribeTopic(sess SESSION, topic string) int {
	ctopic := C.CString(topic)
	rc     := int( C.sol_unsubscribe_topic( C.SOLHANDLE(sess), ctopic ) )
	C.free( unsafe.Pointer(ctopic) )
	return rc
}

func BindQueue(sess SESSION, queue string, fwdmode FWD_MODE, ackmode ACK_MODE) int {
	cqueue := C.CString(queue)
	rc     := int( C.sol_bind_queue( C.SOLHANDLE(sess), cqueue, uint32(fwdmode), uint32(ackmode) ) )
	C.free( unsafe.Pointer(cqueue) )
	return rc
}

func UnbindQueue(sess SESSION, queue string) int {
	cqueue := C.CString(queue)
	rc     := int( C.sol_unbind_queue( C.SOLHANDLE(sess), cqueue ) )
	C.free( unsafe.Pointer(cqueue) )
	return rc
}

func AckMsg(sess SESSION, flow uint64, msgid uint64) int {
	return int( C.sol_ack_msg( C.SOLHANDLE(sess), C.FLOWHANDLE(flow), C.SOLMSGID(msgid) ) )
}

func CacheReq(sess SESSION, cache_name string, topic_sub string , request_id int) int {
	ccache_name := C.CString(cache_name)
	ctopic_sub  := C.CString(topic_sub)
	rc          := int( C.sol_cache_req( C.SOLHANDLE(sess), ccache_name, ctopic_sub, C.int(request_id) ) )
	C.free( unsafe.Pointer(ccache_name) )
	C.free( unsafe.Pointer(ctopic_sub) )
	return rc
}

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

func (this Solace) Connect(host string, vpn string, user string, pass string, appName string, appDesc string, windowsize string) {
	ConnectWithParams(this.sess, host, vpn, user, pass, appName, appDesc, windowsize)
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