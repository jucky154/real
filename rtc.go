/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"C"
	"fmt"
	"time"
	"github.com/tadvi/winc"
	"github.com/gorilla/websocket"
	"github.com/nextzlog/zylo"
	"github.com/recws-org/recws"
	"gopkg.in/go-toast/toast.v1"
	"encoding/json"
	"strconv"
	"sort"
	"strings"
)

var (
	ws = recws.RecConn {
		KeepAliveTimeout: 30 * time.Second,
	}
	url string
	mainWindow *winc.Form
	subWindow *winc.Form
	ls *winc.ListView
	dock *winc.SimpleDock
	dock1 *winc.SimpleDock
	panel[99] *winc.Panel
	ls_section *winc.ListView
	first int
	sub int
	check int
	select_section string
	sections map[string] ([]Station)
	stopCh chan struct{}
    	doneCh chan struct{}
)


type Station struct { 
	CALL   string       `json:"call"`
	SCORE  int          `json:"score"`
	TOTAL  int          `json:"total"`
}

type ByTOTAL []Station
func (a ByTOTAL) Len() int          {return len(a)}
func (a ByTOTAL) Less(i,j int) bool {return a[i].TOTAL>a[j].TOTAL}
func (a ByTOTAL) Swap(i,j int)      {a[i],a[j]=a[j],a[i]}

func notify(msg string) {
	toast := toast.Notification {
		AppID: "ZyLO",
		Title: "ZyLO",
		Message: msg,
	}
	toast.Push()
}

type Item struct {
	T       []string
	checked bool
}

type Atem struct {
	T       []string
}

func (item Item) Text() []string    { return item.T }
func (item *Item) SetText(s string) { item.T[0] = s }

func (item Item) Checked() bool            { return item.checked }
func (item *Item) SetChecked(checked bool) { item.checked = checked }
func (item Item) ImageIndex() int          { return 0 }


//export zlaunch
func zlaunch(cfg string) {
	sub=1
	subWindow = winc.NewForm(nil)
	subWindow.SetSize(400, 300)
	subWindow.SetText("Registration")
	edt := winc.NewEdit(subWindow)
	edt.SetPos(10, 20)
	// Most Controls have default size unless SetSize is called.
	edt.SetText("wss://realtime.allja1.org/agent/")
	btn := winc.NewPushButton(subWindow)
	btn.SetText("Register")
	btn.SetPos(40, 50)
	btn.SetSize(100, 40)
	btn.OnClick().Bind(func(e *winc.Event) {
		url:=edt.ControlBase.Text()
		ws.Dial(url, nil)
		err := ws.GetDialError()
		if err != nil {
			notify(err.Error())
		} else {
			subWindow.Close()
			sub=0
			notify(fmt.Sprintf("successfully connected to %s", url))
			first=1
			check=0
			select_section=""
    			stopCh = make(chan struct{})
   			doneCh = make(chan struct{})
			makemainWindow()
			go onmessage()
		}
		
	})
	subWindow.Center()
	subWindow.Show()
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
func zresult(log uintptr) (total int) {
	total = 0;
	return;
}

const (
	INSERT = 0
	DELETE = 1
)

//export zinsert
func zinsert(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	sendQSO(INSERT, qso)
	notify(fmt.Sprintf("append QSO with %s", qso.GetCall()))
}

//export zdelete
func zdelete(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	sendQSO(DELETE, qso)
	notify(fmt.Sprintf("delete QSO with %s", qso.GetCall()))
}

//export zfinish
func zfinish() {
	if sub==1{
		subWindow.Close()
	}else{	
		close(stopCh)
		<-doneCh
		ws.Close()
		mainWindow.Close()
	}
}

func sendQSO(request byte, qso *zylo.QSO) {
	log := append(*new(zylo.Log), *qso)
	msg := append([]byte{request}, log.Dump(time.Local)...)
	err := ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		notify(err.Error())
	}
}

func onmessage() {
	defer func() { close(doneCh) }()

	for{
		_,data,err := ws.ReadMessage()
		if err == nil {	
			json.Unmarshal(data, &sections)
			reload(sections)
		}
        	select {
        	case <- stopCh:
            		return
		default:
        	}
	}
}

func makemainWindow(){
	// --- Make Window
	mainWindow = winc.NewForm(nil)
	mainWindow.SetSize(800, 600)
	mainWindow.SetText("Ranking")
	dock = winc.NewSimpleDock(mainWindow)
	tabs := winc.NewTabView(mainWindow)
	//check rank
	
	panel[0] = tabs.AddPanel("check rivals")
	edt := winc.NewEdit(panel[0])
	edt.SetPos(10, 20)
	edt.SetText("what is rival's callsign?")

	btn := winc.NewPushButton(panel[0])
	btn.SetText("check!")
	btn.SetPos(40, 50)
	btn.SetSize(100, 40)
	ls_rank := winc.NewListView(panel[0])
	ls_rank.EnableEditLabels(false)
	ls_rank.AddColumn("section", 120)
	ls_rank.AddColumn("rank", 120)
	ls_rank.AddColumn("call sign", 120)
	ls_rank.AddColumn("point", 120)
	ls_rank.AddColumn("score", 120)


	btn.OnClick().Bind(func(e *winc.Event) {
		if sections == nil{
			notify(fmt.Sprintf("none ranking data"))
		} else {
			if check!=0{
				ls_rank.DeleteAllItems()
			}
				
			check=0
			callsign:=edt.ControlBase.Text()
			for section_name,section := range sections {
				sort.Sort(ByTOTAL(section))
				j:=0
				before_score:=-1
				wait_rank:=0
				for _, station := range section {
					if before_score == station.SCORE {
						wait_rank=wait_rank+1
					}else{
						j=j+1+wait_rank
						wait_rank=0
						before_score = station.SCORE
					}
					if strings.Index(station.CALL,callsign)>=0{
						p := &Item{[]string{section_name,strconv.Itoa(j), station.CALL, strconv.Itoa(station.SCORE), strconv.Itoa(station.TOTAL)}, false}
						ls_rank.AddItem(p)
						check=1
					}
				}
			}
			if check==0{
				notify(fmt.Sprintf("your rival doesn't register this contest"))
			}
		}			
	})	
	dock0 := winc.NewSimpleDock(panel[0])
	dock0.Dock(btn, winc.Top)
	dock0.Dock(edt, winc.Top)
	dock0.Dock(ls_rank, winc.Top)

	
	panel[1] = tabs.AddPanel("ranking")
	ls_section = winc.NewListView(panel[1])
	ls_section.EnableEditLabels(false)
	ls_section.AddColumn("select", 200)
	p := &Item{[]string{"ALL"},true}
	ls_section.AddItem(p)

	ls_section.OnClick().Bind(func(e *winc.Event) {
		if sections == nil{
			notify(fmt.Sprintf("none ranking data"))
		} else {
			if ls_section.SelectedCount() > 0{
				item_select:=ls_section.SelectedItem()
				item_select_string:=item_select.Text()
				select_section = item_select_string[0]
				if select_section=="ALL" {
					select_section=""
				}
				reload(sections)
			}
		}			
	})	

	ls = winc.NewListView(panel[1])
	ls.EnableEditLabels(false)
	ls.AddColumn("section", 120)
	ls.AddColumn("rank", 120)
	ls.AddColumn("call sign", 120)
	ls.AddColumn("point", 120)
	ls.AddColumn("score", 120)

	dock1 = winc.NewSimpleDock(panel[1])
	dock1.Dock(ls_section, winc.Left)
	dock1.Dock(ls, winc.Fill)


	// --- Dock(list)
	dock.Dock(tabs, winc.Top)           
	dock.Dock(tabs.Panels(), winc.Fill)
	mainWindow.Center()
	mainWindow.Show()
}
 
func reload(sections map[string] ([]Station)){
	//define section combobox 
	if first == 1{
		for section_name,_ := range sections {
			p := &Item{[]string{section_name}, false}
			ls_section.AddItem(p)
		}
	}
		
	//delete ranking

	if first != 1{
		ls.DeleteAllItems()
	}


	for section_name,section := range sections {
		//reload
		sort.Sort(ByTOTAL(section))
		j:=0
		before_score:=-1
		wait_rank:=0
		for _, station := range section {
			if before_score == station.SCORE {
				wait_rank=wait_rank+1
			}else{
				j=j+1+wait_rank
				wait_rank=0
				before_score = station.SCORE
			}
			if strings.Index(section_name,select_section)>=0{
				p := &Item{[]string{section_name,strconv.Itoa(j), station.CALL, strconv.Itoa(station.SCORE), strconv.Itoa(station.TOTAL)}, false}
				ls.AddItem(p)
			}
		}
		// --- Dock(list and tab)
		dock1.Dock(ls, winc.Fill)
		first=0
	}
}

func main() {}