/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
 */
package main
​
import (
	"C"
	"fmt"
	"time"
	"github.com/tadvi/winc"
	"github.com/gorilla/websocket"
	"github.com/nextzlog/zylo"
	"github.com/recws-org/recws"
	"gopkg.in/go-toast/toast.v1"
)
​
var (
	ws = recws.RecConn {
		KeepAliveTimeout: 30 * time.Second,
	}
	url = "wss://realtime.allja1.org/agent/a6ceda1d-517a-404b-92d5-2b7bfb9bc73b"
	call = "JA1ZLO"
	mainWindow *winc.Form
)
​
func notify(msg string) {
	toast := toast.Notification {
		AppID: "ZyLO",
		Title: "ZyLO",
		Message: msg,
	}
	toast.Push()
}
​
type Item struct {
	T       []string
	checked bool
}
​
func (item Item) Text() []string    { return item.T }
func (item *Item) SetText(s string) { item.T[0] = s }
​
func (item Item) Checked() bool            { return item.checked }
func (item *Item) SetChecked(checked bool) { item.checked = checked }
func (item Item) ImageIndex() int          { return 0 }
​
​
//export zlaunch
func zlaunch(cfg string) {
	ws.Dial(url, nil)
	err := ws.GetDialError()
	if err != nil {
		notify(err.Error())
	} else {
		notify(fmt.Sprintf("successfully connected to %s", url))
		go onmessage(call)
	}
	mainWindow := winc.NewForm(nil)
	dock := winc.NewSimpleDock(mainWindow)
​
	mainWindow.SetSize(700, 600)
	mainWindow.SetText("Controls Demo")
​
	
​
	// --- Tabs
	tabs := winc.NewTabView(mainWindow)
	panel1 := tabs.AddPanel("single op")
	panel2 := tabs.AddPanel("multi op")
	
​
​
	ls := winc.NewListView(panel1)
	ls.EnableEditLabels(false)
	ls.AddColumn("prize", 120)
	ls.AddColumn("callsign", 120)
	ls.AddColumn("point", 120)
	ls.AddColumn("multi", 120)
	ls.AddColumn("total point", 120)
	ls.SetPos(10, 180)
	p1 := &Item{[]string{"First", "JA1ZLO","10","10","100"}, true}
	ls.AddItem(p1)
	p2 := &Item{[]string{"Second", "JA1YWX","9","9","81"}, true}
	ls.AddItem(p2)
	p3 := &Item{[]string{"Third ", "JA1ZGP","1","1","1"}, true}
	ls.AddItem(p3)
​
	ls2 := winc.NewListView(panel2)
	ls2.EnableEditLabels(false)
	ls2.AddColumn("prize", 120)
	ls2.AddColumn("callsign", 120)
	ls2.AddColumn("point", 120)
	ls2.AddColumn("multi", 120)
	ls2.AddColumn("total point", 120)
	ls2.SetPos(10, 180)
	p4 := &Item{[]string{"First", "JA1ZLO","10","10","100"}, true}
	ls2.AddItem(p4)
	p5 := &Item{[]string{"Second", "JA1YWX","9","9","81"}, true}
	ls2.AddItem(p5)
	p6 := &Item{[]string{"Third ", "JA1RL","1","1","1"}, true}
	ls2.AddItem(p6)
	
​
	// --- Dock
	dock1 := winc.NewSimpleDock(panel1)
	dock1.Dock(ls, winc.Fill)
​
	dock2 := winc.NewSimpleDock(panel2)
	dock2.Dock(ls2, winc.Fill)
	tabs.SetCurrent(0)
​
	dock.Dock(tabs, winc.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), winc.Fill) // tab panels dock just below tabs and fill area
​
	mainWindow.Center()
	mainWindow.Show()
}
​
//export zrevise
func zrevise(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	qso.SetMul1(qso.GetRcvd())
}
​
//export zverify
func zverify(ptr uintptr) (score int) {
	score = 1;
	return;
}
​
//export zresult
func zresult(log uintptr) (total int) {
	total = 0;
	return;
}
​
const (
	INSERT = 0
	DELETE = 1
)
​
//export zinsert
func zinsert(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	sendQSO(INSERT, qso)
	notify(fmt.Sprintf("append QSO with %s", qso.GetCall()))
}
​
//export zdelete
func zdelete(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	sendQSO(DELETE, qso)
	notify(fmt.Sprintf("delete QSO with %s", qso.GetCall()))
}
​
//export zfinish
func zfinish() {
	ws.Close()
	mainWindow.Close()
}
​
func sendQSO(request byte, qso *zylo.QSO) {
	log := append(*new(zylo.Log), *qso)
	msg := append([]byte{request}, log.Dump(time.Local)...)
	err := ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		notify(err.Error())
	}
}
​
func onmessage(call string) {
	for {
		_,msg,err := ws.ReadMessage()
		if err == nil {
			notify(string(msg))
		}
	}
}
​
func main() {}