/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
 */
package main

import (
 "C"
 "github.com/nextzlog/zylo"
 "golang.org/x/net/websocket"
 "time"
 "github.com/sqweek/dialog"
)

var(
	url="wss://"
	origin="https://"
)

var ws *websocket.Conn
var wserr interface{}




//export zlaunch
func zlaunch(uintptr) {
	ws, wserr = websocket.Dial(url,"",origin)
	
	if wserr != nil{	
		dialog.Message("%s","Not connected websocket. Plese check websocket.").Info()
		
	}
	if wserr == nil{	
		dialog.Message("%s","You are connecting websocket.").Info()
	}
}

//export zrevise
func zrevise(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	qso.SetMul1(qso.GetRcvd())
}

//export zverify
func zverify(ptr uintptr) (score int) {
	score = 1;
	return;
}

//export zresult
func zresult(qso uintptr) (total int) {
	total = 0;
	return;
}

//export zinsert
func zinsert(ptr uintptr) {
	insert := []byte{0}
	qso:=zylo.ToQSO(ptr)
	log:=new(zylo.Log)
	*log=append(*log,*qso) 
	insert=append(insert,log.Dump(time.Local)...)

	websocket.Message.Send(ws,insert)
	
}

//export zdelete
func zdelete(ptr uintptr) {
	insert := []byte{1}
	qso:=zylo.ToQSO(ptr)
	log:=new(zylo.Log)
	*log=append(*log,*qso) 
	insert=append(insert,log.Dump(time.Local)...)

	websocket.Message.Send(ws,insert)

}

func main() {}
