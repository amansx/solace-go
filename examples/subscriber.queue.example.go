package main

import "fmt"
import "sync";
import "github.com/amansx/solace-go"

const queueName = "myqueue"

func main() {


	var wg sync.WaitGroup
	wg.Add(1)

	s := NewSolace()
	s.Connect("host.docker.internal:55555", "default", "default", "", "1")

	s.NotifyOnMessage(func(msg MessageEvent) {
		fmt.Printf("%+v\n\n", msg)
	})
	
	s.NotifyOnError(func(e ErrorEvent) {
		fmt.Printf("%+v\n", e)
	})

	s.NotifyOnPublisher(func(p PublisherEvent) {
		fmt.Printf("%+v\n", p)
	})
	
	s.NotifyOnConnection(func(e ConnectionEvent) {
		fmt.Printf("%+v\n", e)
	})

	s.Subscribe(queueName)
	
	// s.SubscribeQueue(queueName, ACK_MODE_MANUAL)
	// s.UnsubscribeQueue(queueName)
	// s.Disconnect()
	
	wg.Wait()
}