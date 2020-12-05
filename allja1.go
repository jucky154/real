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

//export zlaunch
func zlaunch(string) {
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
	qso := ToQSO(ptr)
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
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString("Hello")
	qso := ToQSO(ptr)
	file.WriteString(qso.GetCall())
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

//export zdelete
func zdelete(ptr uintptr) {

}

func main() {}
