package main

import "fmt"
import "sync"
import "github.com/amansx/solace-go"

const queueName = "myqueue"
func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	s := NewSolace()
	s.Connect("host.docker.internal:55555", "default", "default", "", "1")

	for i := 0; i < 10; i += 1 {
		size := 10

		b := make([]byte, size)
		for j := 0; j < size; j += 1 {
			b[j] = byte(i)
		}

		id := fmt.Sprintf("%v", i)

		fmt.Println(id)

		s.Publish(
			DESTINATION_TYPE_NONE, queueName, 
			DESTINATION_TYPE_QUEUE, queueName, 
			"JSON", 
			b,
			map[string]interface{}{ 
				"aman": 12, 
				"abc": true, 
				"yada": "aman",
			},
			id,
			i,
		)
	}

	// s.UnsubscribeQueue(queueName)
	// s.Unsubscribe(queueName)
	
	wg.Wait()
	
	s.Disconnect()

}