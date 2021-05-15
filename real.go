/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
	"github.com/nextzlog/zylo"
	"github.com/recws-org/recws"
	"github.com/tadvi/winc"
)

var (
	ws = recws.RecConn{
		KeepAliveTimeout: 30 * time.Second,
	}
	url            string
	geturl         string
	regurl         string
	conurl         string
	path_cfg       string
	flgurl         bool
	mainWindow     *winc.Form
	ls             *winc.ListView
	dock           *winc.SimpleDock
	dock1          *winc.SimpleDock
	panel          [99]*winc.Panel
	ls_section     *winc.ListView
	first          bool
	check          bool
	select_section string
	sections       map[string]([]Station)
	stopCh         chan struct{}
)

type Station struct {
	CALL  string `json:"call"`
	SCORE int    `json:"score"`
	TOTAL int    `json:"total"`
}

type ByTOTAL []Station

func (a ByTOTAL) Len() int           { return len(a) }
func (a ByTOTAL) Less(i, j int) bool { return a[i].TOTAL > a[j].TOTAL }
func (a ByTOTAL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Item struct {
	T       []string
	checked bool
}

type Atem struct {
	T []string
}

func (item Item) Text() []string    { return item.T }
func (item *Item) SetText(s string) { item.T[0] = s }

func (item Item) Checked() bool            { return item.checked }
func (item *Item) SetChecked(checked bool) { item.checked = checked }
func (item Item) ImageIndex() int          { return 0 }

type key struct {
	multinumber string
	band        string
}

var mulmap map[key]int

//go:embed ja1.dat
var ja1list string

func makemap() {
	mulmap = make(map[key]int)
	arr := strings.Fields(ja1list)
	for index, value := range arr {
		if index%2 == 0 {
			for cnt := 0; cnt < 16; cnt++ {
				mulmap[key{value, strconv.Itoa(cnt)}] = 1
			}
		}
	}
}

func conws() {
	url = regurl + geturl
	ws.Dial(url, nil)
	err := ws.GetDialError()
	if err != nil {
		zylo.Notify(err.Error())
	} else {
		zylo.Notify("successfully connected to %s", url)
		stopCh = make(chan struct{})
		go onmessage()
		mainWindow.Show()
	}
}

func checkserver(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get url!")
	r.ParseForm()
	_, tf := r.Form["url_long"]
	if flgurl == false {
		if tf == true {
			for _, v := range r.Form {
				geturl = strings.Join(v, "")
				flgurl = true
				append_cfg()
				conws()
			}
		}
	}
}

//append get url to real.cfg
func append_cfg() {
	//ここでconurlを追記する
	f, _ := os.OpenFile(path_cfg, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()

	fmt.Fprintln(f, "\n"+"conurl  http://localhost:12345/?url_long="+geturl)
}

func makehttp() {
	http.HandleFunc("/", checkserver)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		zylo.Notify(err.Error())
	}
}

func opencfg(path string) {
	path_cfg = path
	cfgdata, _ := ioutil.ReadFile(path)
	cfgarr := strings.Fields(string(cfgdata))
	for index, value := range cfgarr {
		if value == "conurl" {
			conurl = cfgarr[index+1]
		}
		if value == "regurl" {
			regurl = cfgarr[index+1]
		}

	}
}

func zlaunch() {
	makemainWindow()
}

func zattach(name, path string) {
	first = true
	check = false
	select_section = ""
	flgurl = false
	opencfg(path)
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", conurl).Start()
	go makehttp()

	/*
		This is a sample code of adding a QSO to zLog:
		qso := new(zylo.QSO)
		qso.SetCall("JA1FOO")
		qso.SetRcvd("100110")
		qso.Insert()
	*/
}

func zverify(list zylo.Log) (score int) {
	makemap()
	for _, qso := range list {
		call := qso.GetCall()
		rcvd := qso.GetRcvd()
		band := strconv.Itoa(int(qso.Band))
		qso.SetMul1(rcvd)
		if call != "" && mulmap[key{rcvd, band}] > 0 {
			score = 1
			if mulmap[key{rcvd, band}] == 1 {
				qso.SetNewMul1(true)
			}
			if mulmap[key{rcvd, band}] > 1 {
				qso.SetNewMul1(false)
			}
			mulmap[key{rcvd, band}] = mulmap[key{rcvd, band}] + 1
		}
	}
	return
}

func zupdate(list zylo.Log) (total int) {
	calls := mapset.NewSet()
	mults := mapset.NewSet()
	for _, qso := range list {
		call := qso.GetCall()
		mul1 := qso.GetMul1()
		new1 := !mults.Contains(mul1)
		qso.SetNewMul1(new1)
		calls.Add(call)
		mults.Add(mul1)
	}
	score := calls.Cardinality()
	multi := mults.Cardinality()
	total = score * multi
	return
}

const (
	INSERT = 0
	DELETE = 1
)

func zinsert(list zylo.Log) {
	for _, qso := range list {
		sendQSO(INSERT, qso)
		zylo.Notify("append QSO with %s", qso.GetCall())
	}
}

func zdelete(list zylo.Log) {
	for _, qso := range list {
		sendQSO(DELETE, qso)
		zylo.Notify("delete QSO with %s", qso.GetCall())
	}
}

func zkpress(key int, source string) (block bool) {
	block = false
	return
}

func zfclick(btn int, source string) (block bool) {
	block = false
	return
}

func zdetach() {
	if ws.IsConnected() {
		close(stopCh)
		ws.Close()
	}
	mainWindow.Close()
}

func zfinish() {}

func sendQSO(request byte, qso zylo.QSO) {
	log := append(*new(zylo.Log), qso)
	msg := append([]byte{request}, log.Dump(time.Local)...)
	err := ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		zylo.Notify(err.Error())
	}
}

func onmessage() {
	for {
		select {
		case <-stopCh:
			zylo.Notify("real.dll stop routine")
			return
		default:
			_, data, err := ws.ReadMessage()
			if err == nil {
				json.Unmarshal(data, &sections)
				reload(sections)
			}
		}
	}
}

func makemainWindow() {
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
		if sections == nil {
			zylo.Notify("none ranking data")
		} else {
			if check {
				ls_rank.DeleteAllItems()
			}

			check = false
			callsign := edt.ControlBase.Text()
			for section_name, section := range sections {
				sort.Sort(ByTOTAL(section))
				j := 0
				before_score := -1
				wait_rank := 0
				for _, station := range section {
					if before_score == station.TOTAL {
						wait_rank = wait_rank + 1
					} else {
						j = j + 1 + wait_rank
						wait_rank = 0
						before_score = station.TOTAL
					}
					if strings.Index(station.CALL, callsign) >= 0 {
						p := &Item{[]string{section_name, strconv.Itoa(j), station.CALL, strconv.Itoa(station.SCORE), strconv.Itoa(station.TOTAL)}, false}
						ls_rank.AddItem(p)
						check = true
					}
				}
			}
			if !check {
				zylo.Notify("your rival doesn't register this contest")
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
	p := &Item{[]string{"ALL"}, true}
	ls_section.AddItem(p)

	ls_section.OnClick().Bind(func(e *winc.Event) {
		if ls_section.SelectedCount() == 1 {
			if sections == nil {
				zylo.Notify("none ranking data")
			} else {
				item_select := ls_section.SelectedItem()
				item_select_string := item_select.Text()
				select_section = item_select_string[0]
				if select_section == "ALL" {
					select_section = ""
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
}

func reload(sections map[string]([]Station)) {
	//define section combobox
	if first {
		for section_name, _ := range sections {
			p := &Item{[]string{section_name}, false}
			ls_section.AddItem(p)
		}
	}

	//delete ranking

	if !first {
		ls.DeleteAllItems()
	}

	for section_name, section := range sections {
		//reload
		sort.Sort(ByTOTAL(section))
		j := 0
		before_score := -1
		wait_rank := 0
		for _, station := range section {
			if before_score == station.TOTAL {
				wait_rank = wait_rank + 1
			} else {
				j = j + 1 + wait_rank
				wait_rank = 0
				before_score = station.TOTAL
			}
			if strings.Index(section_name, select_section) >= 0 {
				p := &Item{[]string{section_name, strconv.Itoa(j), station.CALL, strconv.Itoa(station.SCORE), strconv.Itoa(station.TOTAL)}, false}
				ls.AddItem(p)
			}
		}
		// --- Dock(list and tab)
		dock1.Dock(ls, winc.Fill)
		first = false
	}
}

func main() {}
