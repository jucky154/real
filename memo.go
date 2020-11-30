package main

import (
		"log"
		"time"
		"strings"
		"golang.org/x/net/websocket"
)

	var (
		origin = "http://localhost:9000/"
		url = "ws://localhost:9000/flash"
)

type EchoMsg struct{
	zlog unit64
}

func main(){
	ws, err := websocket.Dial("url", "", "origin");
 	if err != nil {
		panic("Dial: " + err.String())
	}
	if _, err := ws.Write(zlog); err != nil {
		panic("Write: " + err.String())
	}
	var msg = make([]byte, 512);
	if n, err := ws.Read(msg); err != nil {
		panic("Read: " + err.String())
	}
}