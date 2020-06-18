package main

import "fmt";
import "time";
import "plugin";

func main() {
	if so, err := plugin.Open("./solace.goplugin"); err == nil {

		if gosolace, err := so.Lookup("GoSolace"); err == nil {

			if gosol, ok := gosolace.(func() interface{}); ok {

				solace := gosol().(Solace)
				
				solace.Connect("google.com:55555", "default", "default", "", "1")
				
				solace.Subscribe("aman")
				
				solace.Publish("aman", []byte("Hello World"))
				
				solace.Unsubscribe("aman")
				
				time.Sleep(10 * time.Second)

				return

			}

			fmt.Println(err)
		}

		fmt.Println(err)
	}

}