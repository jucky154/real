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
 "os"
)

var(
	url="wss://realtime.allja1.org/agent/JA1ZLO/a6ceda1d-517a-404b-92d5-2b7bfb9bc73b"
	origin="https://realtime.allja1.org/agent/JA1ZLO/"
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
	go recvMsg()
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
	sendQso(insert,ptr)
}

//export zdelete
func zdelete(ptr uintptr) {
	insert := []byte{1}
	sendQso(insert,ptr)
}

//export zfinish
func zfinish(){
	_=ws.Close()
	
}

//export sendQso
func sendQso(insert []byte,ptr uintptr){
	qso:=zylo.ToQSO(ptr)
	log:=new(zylo.Log)
	*log=append(*log,*qso) 
	insert=append(insert,log.Dump(time.Local)...)	
	
	err := websocket.Message.Send(ws,insert)

	time.Sleep(time.Second*2)

	if err != nil{	
		file,_ := os.Create("err.ZLO")
		defer file.Close()
		file.Write(log.Dump(time.Local))

		dialog.Message("%s","Not connected websocket. Plese check websocket.make zlo file to send again.").Info()
	}
}

//export recvMsg
func recvMsg(){
	for{
		Msg:="hoge"
		err := websocket.Message.Receive(ws,Msg)
		time.Sleep(time.Second*1)
		if Msg != "hoge"{
			dialog.Message("%s",Msg).Info()
		}

		if err != nil{
			dialog.Message("%s","Not connected websocket. Plese check websocket.Restart zLog ").Info()
		}
	}
}

func main() {}