package main

type Solace interface {
	Init(func(interface{}))
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
}