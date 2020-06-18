package main

type Solace interface {
	Init()
	Connect(string, string, string, string, string)
	Disconnect()
	Subscribe(string)
	Unsubscribe(string)
	Publish(string, []byte)
	PublishDirect(string, []byte)
	PublishToQueue(string, []byte)
	BindQueue(string)
	BindQueueManualAck(string)
	UnbindQueue(string)
	Ack(uint64, uint64)
}