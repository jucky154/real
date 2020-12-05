/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
 */
package main

import (
	"C"
	"os"
	"time"
	"github.com/nextzlog/zylo"
)

//export zlaunch
func zlaunch(cfg string) {
}

//export zrevise
func zrevise(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	qso.SetMul1(qso.GetRcvd())
	qso.SetMul2("")
}

//export zverify
func zverify(ptr uintptr) (score int) {
	return 114514
}

//export zresult
func zresult(log uintptr) (total int) {
	return 364364
}

//export zinsert
func zinsert(ptr uintptr) {
	qso := zylo.ToQSO(ptr)
	log := new(zylo.Log)
	*log = append(*log, *qso)
	file, _ := os.Create("insert.zlo")
	defer file.Close()
	file.Write(log.Dump(time.Local))
}

//export zdelete
func zdelete(ptr uintptr) {}

func main() {}
