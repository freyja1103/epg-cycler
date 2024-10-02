package main

import "fmt"

func FileFoundLog(s string) (int, error) {
	return DebugLog("File found", s)
}

func SCLog(title, subtitle, number string) (int, error) {
	return fmt.Printf("TS info from EPG:\nTitle:	%v\nSubtitle:	%v	Episode:	%v\n", title, subtitle, number)
}

func DebugLog(msg string, s ...string) (int, error) {
	return fmt.Printf("%v:	%v\n", msg, s)
}
