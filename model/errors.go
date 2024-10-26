package model

import (
	"errors"
	"log"
	"runtime"
	"strconv"
	"strings"
)

var (
	UnSupoortedOSError       = errors.New("Current OS is not supported: " + runtime.GOOS)
	UnSupoortedCharCodeError = errors.New("Unsupported character code: Garbled text might be occurring. \nIf you are using a .bat file in the command prompt, please add \"chcp 65001\" at the beginning to load it as UTF-8 encoding.\n")
	WarnProgramName          = errors.New("The format of the program name is not supported. The name of the created directory may differ from the actual program name.")
)

func FileUnsatisfiedError(i int) error {
	return errors.New("File unsatisfied error: target files are required but " + strconv.Itoa(i))
}

func Errorlog(e error) {
	log.Printf("Error Occured:	%v", e)
}

func SubtitleNotFoundError(s string) error {
	return errors.New("Not found subtitle :" + s)
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
		log.Println("Loaded title:	", s)
		return UnSupoortedCharCodeError
	}
	return nil
}
