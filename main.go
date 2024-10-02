package main

import (
	"encoding/xml"
	"flag"
	"fmt"
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
	args[1] = "originpath"
	args[2] = "title"
	args[3] = "basename"
	args[4] = "number"
	args[5] = "process"
	args[6] = "ip"
	save_path := flag.String(args[0], "", "save video path")
	origin_path := flag.String(args[1], "", "origin video path")
	title := flag.String(args[2], "", "a program's name")
	basename := flag.String(args[3], "", "filename without ext")
	number := flag.String(args[4], "number", "episode number")
	procPreventShutdown := flag.String(args[5], "", "process that prevent shutdown")
	APIURL := flag.String(args[6], "localhost:5510", "EpgTimer's HTTP server, IP:port")
	flag.Parse()

	if err := CheckArg(*title, args); err != nil {
		Errorlog(err)
		return
	}

	SCLog(*title, *basename, *number)
	err := dirTidy(save_path, origin_path, title, basename)
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

func dirTidy(save_path, origin_path, title, basename *string) error {
	var files = []string{*origin_path, *origin_path + ".err", *origin_path + ".program.txt"}
	var half_title string
	if isNotValidFilename(*title) {
		err := OperateFile(*save_path, *origin_path, *title, files)
		if err != nil {
			return err
		}
	} else {
		half_title = width.Fold.String(*title)
		err := OperateFile(*save_path, *origin_path, half_title, files)
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

func isNotValidFilename(s string) bool {
	invalidChars := `<>:"/\|?*`
	for _, char := range invalidChars {
		if strings.ContainsRune(s, char) {
			return true
		}
	}
	return false
}

func OperateFile(save_path, origin_path, title string, files []string) error {
	err := os.Mkdir(title, 0755)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	for _, file := range files {
		err := os.Rename((filepath.Join(filepath.Dir(origin_path), file)), filepath.Join(save_path, title, width.Fold.String(file)))
		if err != nil {
			return err
		}
	}
	return nil
}
