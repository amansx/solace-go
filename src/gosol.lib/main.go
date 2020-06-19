package main

import "gosol_t"

var singleton *Solace

func GoSolace(onMessage gosol_t.MessageHandler, onError gosol_t.ErrorHandler, onConnection gosol_t.ConnectionEventHandler, onPublisherEvent gosol_t.PublisherEventHandler) gosol_t.Solace {
	if (singleton == nil) {
		singleton = &Solace{}
		singleton.Init(onMessage, onError, onConnection, onPublisherEvent)
	}
	return singleton
}