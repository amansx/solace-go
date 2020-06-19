package main

import "fmt";
import "plugin";
import "gosol_t";

func onMessage(e gosol_t.MessageEvent) {
	fmt.Printf("%+v\n", e)
}

func onError(e gosol_t.ErrorEvent) {
	fmt.Printf("%+v\n", e)
}

func onConnectionEvent(e gosol_t.ConnectionEvent) {
	fmt.Printf("%+v\n", e)
}

func onPublisherEvent(e gosol_t.PublisherEvent) {
	fmt.Printf("%+v\n", e)
}

func main() {
	if so, err := plugin.Open(pluginPath("solace.goplugin")); err == nil {
		if gosolace, err := so.Lookup("GoSolace"); err == nil {
			if initSolace, ok := gosolace.(func (gosol_t.MessageHandler, gosol_t.ErrorHandler, gosol_t.ConnectionEventHandler, gosol_t.PublisherEventHandler) gosol_t.Solace); ok {

				solace := initSolace(onMessage, onError, onConnectionEvent, onPublisherEvent)

				solace.Connect("host.docker.internal:55555", "default", "default", "", "1")
				solace.Subscribe("aman")
				solace.Publish("aman", []byte("Hello World"))
				solace.Unsubscribe("aman")

				return

			}

		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println(err)
	}

}