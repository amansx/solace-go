package main

var singleton Solace

type driverinitializer string
func (this driverinitializer) Load() {
	singleton = Solace{}
	singleton.Init()
}

var GoSolace driverinitializer