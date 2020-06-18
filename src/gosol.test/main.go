package main

import "fmt";
import "time";
import "plugin";

func main() {
	if so, err := plugin.Open(execPath() + "/solace.goplugin"); err == nil {

		if gosolace, err := so.Lookup("GoSolace"); err == nil {

			if gosol, ok := gosolace.(func() interface{}); ok {

				solace := gosol().(Solace)
				
				solace.Connect("host.docker.internal:55555", "default", "default", "", "1")
				
				solace.Subscribe("aman")
				
				solace.Publish("aman", []byte("Hello World"))
				
				solace.Unsubscribe("aman")
				
				time.Sleep(2 * time.Second)

				return

			}

		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println(err)
	}

}