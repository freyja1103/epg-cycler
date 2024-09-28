package main

import "fmt"

func FileFoundLog(s string) (int, error) {
	return DebugLog("File found", s)
}

func DebugLog(msg string, s ...string) (int, error) {
	return fmt.Printf("%v:	%v\n", msg, s)
}
