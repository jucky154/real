/*
 implements the real time contest client software.
 Copyright (C) 2020 JA1ZLO.
 */
package main

import "C"
import "github.com/nextzlog/zylo"

//export zlaunch
func zlaunch(warning zylo.Warning) {

}

//export zrevise
func zrevise(qso *zylo.QSO) {

}

//export zverify
func zverify(qso *zylo.QSO) (score int) {
	score = 0;
	return;
}

//export zresult
func zresult(qso *zylo.Log) (total int) {
	total = 0;
	return;
}

//export zinsert
func zinsert(qso *zylo.QSO) {

}

//export zdelete
func zdelete(qso *zylo.QSO) {

}

func main() {}
