package main

import "os";
import "fmt";
import "plugin";

type GoSolace interface {
	Init()
}

func main() {
	if so, err := plugin.Open("./solace.plugin"); err == nil {
		if gosolace, err := so.Lookup("GoSolace"); err == nil {
			if solace, ok := gosolace.(GoSolace); ok {
				solace.Load()
				return
			} 
			fmt.Println(err)
		}
		fmt.Println(err)
	}
}