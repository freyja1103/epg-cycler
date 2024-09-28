package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/text/width"
)

func main() {
	args := make([]string, 6)
	args[0] = "srcpath"
	args[1] = "title"
	args[2] = "subtitle"
	args[3] = "number"
	args[4] = "process"
	args[5] = "ip"

	srcpath := flag.String(args[0], "", "save video path")
	SCtitle := flag.String(args[1], "title", "Anime's name")
	SCsubtitle := flag.String(args[2], "subtitle", "episode title")
	SCnumber := flag.String(args[3], "number", "episode number")
	procPreventShutdown := flag.String(args[4], "", "process that prevent shutdown")
	APIURL := flag.String(args[5], "localhost:5510", "EpgTimer's HTTP server, IP:port")
	flag.Parse()

	if err := CheckArg(*SCtitle, args); err != nil {
		Errorlog(err)
		return
	}

	err := dirTidy(srcpath, SCtitle, SCsubtitle, SCnumber)
	if err != nil {
		fmt.Println("Error occured:", err)
		return
	}

	url := "http://" + *APIURL + "/api/EnumReserveInfo"
	body, err := APIReq2Body(url)
	if err != nil {
		Errorlog(err)
	}

	var entry Entry
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		Errorlog(err)
	}

	hasReserve, err := HasRemainReserve(&entry)
	if hasReserve {
		if err != nil {
			Errorlog(err)
		}
		return
	}

	isExec, err := NoShutdownTrigger(*procPreventShutdown)
	if err != nil {
		Errorlog(err)
		return
	}
	if isExec {
		return
	}

	// shutdown
	if ExecShutdown() != nil {
		Errorlog(err)
	}
}

func dirTidy(srcpath, SCtitle, SCsubtitle, SCnumber *string) error {
	var files []string
	var err error
	half_title := width.Fold.String(*SCtitle)
	err = filepath.WalkDir(*srcpath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		FileFoundLog(path)
		if !strings.Contains(path, *SCtitle) || info.IsDir() {
			return nil
		}

		if !(filepath.Dir(path) == filepath.Join(*srcpath, *SCtitle)) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return FileUnsatisfiedError(len(files))
	}

	err = os.Mkdir(half_title, 0755)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	for _, file := range files {
		err := os.Rename(file, filepath.Join(filepath.Dir(file), half_title, width.Fold.String(filepath.Base(file))))
		if err != nil {
			return err
		}
	}

	return nil
}

func NoShutdownTrigger(targetProcess string) (bool, error) {
	if targetProcess == "" {
		return false, nil
	}

	processes, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range processes {
		name, err := p.Name()
		if err == nil && name == targetProcess {
			return true, nil
		}
		continue
	}

	return false, nil
}

func ExecShutdown() error {
	if runtime.GOOS == "windows" {
		exec := exec.Command(
			"C:\\Windows\\System32\\shutdown.exe",
			"/s /t 60 /f /c 'shutdown by epg-cycler after 60s'",
		)
		if exec.Err != nil {
			return exec.Err
		}
		return nil
	}
	if runtime.GOOS == "darwin" {
		exec := exec.Command("shutdown", "-h now")
		if exec.Err != nil {
			return exec.Err
		}
		return nil
	}

	return UnSupoortedOSError()
}
