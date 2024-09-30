package main

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func UnSupoortedOSError() error {
	return errors.New("This OS is not supported: " + runtime.GOOS)
}

func FileUnsatisfiedError(i int) error {
	return errors.New("File unsatisfied error: target files are required but " + strconv.Itoa(i))
}

func UnSupoortedCharCodeError(s string) error {
	return errors.New("Unsupported character code: Garbled text might be occurring. \nIf you are using a .bat file in the command prompt, please add \"chcp 65001\" at the beginning to load it as UTF-8 encoding.\n")
}

func Errorlog(e error) (int, error) {
	return fmt.Printf("Error Occured:	%v", e)
}

func IsGrabledText(s string, args []string) bool {
	for _, arg := range args {
		if strings.Contains(s, "-"+arg+"=") {
			return true
		}
	}
	return false
}

func CheckArg(s string, args []string) error {
	if IsGrabledText(s, args) {
		fmt.Println("Loaded title:	", s)
		return UnSupoortedCharCodeError(s)
	}
	return nil
}
