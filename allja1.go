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

func pointer_to_qso(ptr uintptr) *zylo.QSO {
	return (*zylo.QSO)(unsafe.Pointer(ptr))
}

//export zlaunch
func zlaunch(uintptr) {
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString("Hello")
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

//export zrevise
func zrevise(ptr uintptr) {

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
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString("Hello")
	qso := pointer_to_qso(ptr)
	file.WriteString(qso.GetCall())
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

//export zdelete
func zdelete(ptr uintptr) {

}

func main() {}
