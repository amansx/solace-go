package main

import "C"
import "fmt"
import "gosol"

func on_err(sess gosol.SESSION, err *gosol.ErrEvent) {
	fmt.Println("\nERROR EVENT:")
	fmt.Println("\tFNName: ",     err.FNName)
	fmt.Println("\tRetCcode: ",   err.RetCode)
	fmt.Println("\tRCString: ",   err.RCStr)
	fmt.Println("\tSubCode: ",    err.SubCode)
	fmt.Println("\tSCString: ",   err.SCStr)
	fmt.Println("\tRespCode: ",   err.RespCode)
	fmt.Println("\tErr String: ", err.ErrStr)
}

func on_msg(sess gosol.SESSION, msg *gosol.MsgEvent) {
	fmt.Println("\nMESSAGE EVENT:")

	cstr    := (*C.char)(msg.Buffer)
	payload := C.GoStringN(cstr, C.int(msg.BufLen))

	fmt.Println("\tDestination: ", msg.Destination)
	fmt.Println("\tBuffer: ",      payload)
	fmt.Println("\tBufLen: ",      msg.BufLen)
	fmt.Println("\tMsgId: ",       msg.MsgId)
	fmt.Println("\tRedelivered: ", msg.Redelivered)
	fmt.Println("\tDiscard: ",     msg.Discard)
}

func on_con(sess gosol.SESSION, con *gosol.ConEvent) {
	fmt.Println("\nCONNECTIVITY EVENT:")
	fmt.Println("\tType: ", connectionTypeToString(con.Type))
}

func on_pub(sess gosol.SESSION, pub *gosol.PubEvent) {
	fmt.Println("\nPUBLISHER EVENT:")
	fmt.Println("\tType: ", publishedEventTypeToString(pub.Type))
	correlation := *(*C.int)(pub.CorrelationData)
	fmt.Println("\tCorrelationData: ", correlation)
}
