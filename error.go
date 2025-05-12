package main

import (
	"context"
	"errors"
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/freyja1103/epg-cycler/logging"
)

func UnSupoortedOSError() error {
	return errors.New("this OS is not supported: " + runtime.GOOS)
}

func FileUnsatisfiedError(i int) error {
	return errors.New("file unsatisfied error: target files are required but " + strconv.Itoa(i))
}

func UnSupoortedCharCodeError(s string) error {
	return errors.New("unsupported character code: Garbled text might be occurring. \nIf you are using a .bat file in the command prompt, please add \"chcp 65001\" at the beginning to load it as UTF-8 encoding.\n")
}

func SubtitleNotFoundError(s string) error {
	return errors.New("not found subtitle :" + s)
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
		logging.InfoAttrs(context.Background(), slog.LevelInfo, "grabled text detected", slog.String("grabled text", s))
		return UnSupoortedCharCodeError(s)
	}
	return nil
}
