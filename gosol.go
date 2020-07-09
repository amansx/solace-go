package main

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
import cptr "github.com/mattn/go-pointer"

type SESSION uint64

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

func ConnectWithParams(sess SESSION, host string, vpn string, user string, pass string, windowsize string) int {
	// C.CString allocates new memory on the heap, so need to cleanup afterwards
	chost := C.CString(host)
	cvpn  := C.CString(vpn)
	cuser := C.CString(user)
	cpass := C.CString(pass)
	cws   := C.CString(windowsize)

	rc    := int( C.sol_connect_with_params( C.SOLHANDLE(sess), chost, cvpn, cuser, cpass, cws ) )

	C.free(unsafe.Pointer(chost))
	C.free(unsafe.Pointer(cvpn))
	C.free(unsafe.Pointer(cuser))
	C.free(unsafe.Pointer(cpass))
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


