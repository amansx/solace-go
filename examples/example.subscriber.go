package main

import "fmt"
import "sync";
import solace "github.com/amansx/solace-go"

const queueName = "myqueue"

func main() {


	var wg sync.WaitGroup
	wg.Add(1)

	s := solace.NewSolace()
	s.Connect("host.docker.internal:55555", "default", "default", "", "1", "", "")

	s.NotifyOnMessage(func(msg solace.MessageEvent) {
		fmt.Printf("%+v\n\n", msg)
	})
	
	s.NotifyOnError(func(e solace.ErrorEvent) {
		fmt.Printf("%+v\n", e)
	})

	s.NotifyOnPublisher(func(p solace.PublisherEvent) {
		fmt.Printf("%+v\n", p)
	})
	
	s.NotifyOnConnection(func(e solace.ConnectionEvent) {
		fmt.Printf("%+v\n", e)
	})

	s.Subscribe(queueName)
	
	// s.SubscribeQueue(queueName, ACK_MODE_MANUAL)
	// s.UnsubscribeQueue(queueName)
	// s.Disconnect()
	
	wg.Wait()
}