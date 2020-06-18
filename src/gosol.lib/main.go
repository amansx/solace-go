package main

var singleton *Solace

func GoSolace() interface{} {
	singleton = &Solace{}
	singleton.Init()
	return singleton
}