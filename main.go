package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/text/width"
)



func main() {
	srcpath := flag.String("srcpath", "F:\\Videos\\TVRec", "save video path")
	SCtitle := flag.String("title", "title", "Anime's name")
	SCsubtitle := flag.String("subtitle", "subtitle", "episode title")
	SCnumber := flag.String("number", "number", "episode number")
	procPreventShutdown := flag.String("process", "", "process that prevent shutdown")
	APIURL := flag.String("url", "localhost:5510", "EDCB IP:port")
	flag.Parse()

	err := dirTidy(srcpath, SCtitle, SCsubtitle, SCnumber)
	if err != nil {
		fmt.Println("Error occured:", err)
		return
	}
	
	url := "http://" + *APIURL + "/api/EnumReserveInfo"
	body, err := APIReq2Body(url)
	if err != nil {
		fmt.Println("Error Occured:", err)
	}
	
	var entry Entry
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		fmt.Println(err)
	}

	hasReserve, err := HasRemainReserve(&entry)
	if hasReserve {
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	isExec, err := NoShutdownTrigger(*procPreventShutdown);
	if err != nil {
		fmt.Println(err)
		return
	}
	if isExec {
		return
	}

	// shutdown
	if ExecShutdown() != nil {
		fmt.Println(err)
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

		if (!strings.Contains(path, *SCtitle) || info.IsDir()) {
			return nil
		}
		if !(filepath.Dir(path) == filepath.Join(*srcpath, half_title)) {
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

func UnSupoortedOSError() error {
	return errors.New("This OS is not supported: " + runtime.GOOS)
}

func FileUnsatisfiedError(i int) error {
	return errors.New("File unsatisfied error: target files are required but " + strconv.Itoa(i))
}