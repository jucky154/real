/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
 */
package main

import (
	"C"
	"fmt"
	"log"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/nextzlog/zylo"
	
)



//export zlaunch
func zlaunch(cfg string) {
	go window()
	
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


//export zinsert
func zinsert(ptr uintptr) {
}

//export zdelete
func zdelete(ptr uintptr) {
}

//export zfinish
func zfinish() {
}

func window(){
// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:           "Test",
		BaseDirectoryPath: "example",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("https://realtime.allja1.org/lists", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// Blocking pattern
	a.Wait()
}


func main() {}