package main

import (
	"log"
)

func FileFoundLog(s string) {
	DebugLog("File found", s)
}

func SrcLog(title, basename, number string) {
	log.Printf("TS info from EPG: Title:	%v Basename: %v	Episode: %v\n", title, basename, number)
}

func DebugLog(msg string, s ...string) {
	if s == nil {
		log.Println(msg)
		return
	}
	log.Printf("%v:	%v\n", msg, s)
}
