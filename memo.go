/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
 */
package main

import "C"
import "github.com/nextzlog/zylo"
import "log"
import "os"
import "unsafe"
import "golang.org/x/net/websocket"
import "time"

var(
	url="wss://"
	origin="https://"
)

func pointer_to_qso(ptr uintptr) *zylo.QSO {
	return (*zylo.QSO)(unsafe.Pointer(ptr))
}

//export zlaunch
func zlaunch(uintptr) {
	ws, err := websocket.Dial(url,"",origin)
	_=ws.Close()
	
	if err != nil{	
		file, err := os.Create("output.txt")
		if err != nil {
			log.Fatal(err)
		}

		file.WriteString("No")

		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}
	if err == nil{	
		file, err := os.Create("output.txt")
		if err != nil {
			log.Fatal(err)
		}

		file.WriteString("Yes")

		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

//export zrevise
func zrevise(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	qso.SetMul1(qso.GetRcvd())
}

//export zverify
func zverify(ptr uintptr) (score int) {
	score = 0;
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

	ws, _ := websocket.Dial(url,"",origin)
	websocket.Message.Send(ws,insert)
	_=ws.Close()
	
}

//export zdelete
func zdelete(ptr uintptr) {

}

func main() {}
