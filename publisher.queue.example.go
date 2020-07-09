// +build examples

package main

// import "sync";
import "fmt";
import "os"
import "path"
import "plugin";
import "path/filepath"
import "internal/solace"

func pluginPath(pluginName string) string {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    return path.Join(exPath, pluginName)
}

func onMessage(e solace.MessageEvent) {
	fmt.Printf("%+v\n", e)
	fmt.Println(e.Buffer)
	fmt.Println(e.UserProperties)
}

func onError(e solace.ErrorEvent) {
	fmt.Printf("%+v\n", e)
}

func onConnectionEvent(e solace.ConnectionEvent) {
	fmt.Printf("%+v\n", e)
}

func onPublisherEvent(e solace.PublisherEvent) {
	fmt.Printf("%+v\n", e)
}

const queueName = "myqueue"

func main() {

	// var wg sync.WaitGroup
	// wg.Add(1)

		// Load plugin
		if so, err := plugin.Open(pluginPath("solace.linux.amd64.gopl")); err == nil {

			// Lookup plugin entry function
			if gosolace, err := so.Lookup("InitSolaceWithCallback"); err == nil {

					// Cast Entry function
					if initSolace, ok := gosolace.(*solace.InitSolaceWithCallback); ok {

						s := (*initSolace)(onMessage, onError, onConnectionEvent, onPublisherEvent)
						s.Connect("host.docker.internal:55555", "default", "default", "", "1")


						for i := 0; i < 10000; i += 1 {
							b := make([]byte, 100000)
							fmt.Printf("%+v\n", b)
							s.Publish(
								// solace.DESTINATION_TYPE_QUEUE, 
								solace.DESTINATION_TYPE_NONE, queueName, 
								solace.DESTINATION_TYPE_QUEUE, queueName, 
								"JSON", 
								b,
								map[string]interface{}{ 
									"aman": 12, 
									"abc": true, 
									"yada": "aman",
								},
								"12",
								12,
							)
						}

						// s.UnsubscribeQueue(queueName)
						// s.Unsubscribe(queueName)
						s.Disconnect()

					}
				
			} else {
				fmt.Println(err)
			}

		} else {
			fmt.Println(err)
		}


	// wg.Wait()
}