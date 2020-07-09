// +build examples

package main

import "sync";
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
	fmt.Printf("%v\n", e.Buffer)
	// fmt.Printf("%+v\n", e.BufferLen)
	// fmt.Println(string(e.Buffer))
	// fmt.Println(e.UserProperties)
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

	var wg sync.WaitGroup
	wg.Add(1)

		// Load plugin
		if so, err := plugin.Open(pluginPath("solace.linux.amd64.gopl")); err == nil {
			// Lookup plugin entry function
			if gosolace, err := so.Lookup("InitSolaceWithCallback"); err == nil {
				// Cast Entry function
				if initSolace, ok := gosolace.(*solace.InitSolaceWithCallback); ok {

					go func(){
						s := (*initSolace)(onMessage, onError, onConnectionEvent, onPublisherEvent)
						s.Connect("host.docker.internal:55555", "default", "default", "", "1")
						s.Subscribe(queueName)
						// s.SubscribeQueue(queueName, solace.ACK_MODE_MANUAL)
						// s.UnsubscribeQueue(queueName)
						// s.Disconnect()
					}()

				}

			} else {
				fmt.Println(err)
			}

		} else {
			fmt.Println(err)
		}


	wg.Wait()
}