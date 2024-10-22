package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/text/width"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var tp targetProcesses

	args := make([]string, 8)
	args[0] = "srcpath"
	args[1] = "originpath"
	args[2] = "title"
	args[3] = "basename"
	args[4] = "number"
	args[5] = "process"
	args[6] = "ip"
	args[7] = "all-tidy"
	s_path := flag.String(args[0], "", "save video path")
	o_path := flag.String(args[1], "", "origin video path")
	title := flag.String(args[2], "", "a program's name")
	basename := flag.String(args[3], "", "filename without ext")
	number := flag.String(args[4], "number", "episode number")
	flag.Var(&tp, "process", "process that prevent shutdown")
	APIURL := flag.String(args[6], "localhost:5510", "EpgTimer's HTTP server, IP:port")
	all_tidy_mode := flag.Bool(args[7], false, "The mode for sorting all the recording files in the directory at once")
	flag.Parse()

	save_path := strings.ReplaceAll(filepath.ToSlash(*s_path), "//", "/")
	origin_path := strings.ReplaceAll(filepath.ToSlash(*o_path), "//", "/")

	// 一括整理モード
	if *all_tidy_mode {
		if err := TidyAllFiles(save_path); err != nil {
			Errorlog(err)
			return
		}
		return
	}

	if !(*all_tidy_mode) {
		if err := CheckArg(*title, args); err != nil {
			Errorlog(err)
			return
		}

		SrcLog(*title, *basename, *number)
		err := dirTidy(save_path, origin_path, title, basename)
		if err != nil {
			Errorlog(err)
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

		for _, p := range tp {
			isExec, err := NoShutdownTrigger(p)
			if err != nil {
				Errorlog(err)
				return
			}
			if isExec {
				return
			}
		}

		// shutdown
		if ExecShutdown() != nil {
			Errorlog(err)
		}
	}

}

func dirTidy(save_path, origin_path string, title, basename *string) error {
	// basename -> $FileName$
	var (
		ts_filename      string   = *basename + ".ts"
		fold_ts_filename string   = width.Fold.String(*basename + ".ts")
		half_title       string   = width.Fold.String(*title)
		files            []string = []string{ts_filename, ts_filename + ".err", ts_filename + ".program.txt"}
		converted_files  []string
		subtitle         string
		program_name     string
		err              error
	)
	origin_program_name, ep_string := GetProgramName(*basename)
	isInvalid, _ := isInvalidName(*basename)
	if isInvalid {
		DebugLog("Will be invalid filename, no convert fold style")
		subtitle, err = GetSubtitle(fold_ts_filename)
		if err != nil {
			return err
		}

		isInvalidSubtitle, _ := isInvalidName(subtitle)
		isInvalidProgramName, _ := isInvalidName(origin_program_name)
		if !isInvalidProgramName && isInvalidSubtitle {
			// only subtitle is invalid
			conv_ts_filename := ConcatFilename(*basename, half_title, width.Widen.String(subtitle), ep_string, ".ts")
			converted_files = []string{conv_ts_filename, conv_ts_filename + ".err", conv_ts_filename + ".program.txt"}
			program_name = origin_program_name

		} else {
			half_title = *title
			if strings.HasSuffix(origin_program_name, " ") {
				origin_program_name = origin_program_name[:len(origin_program_name)-1]
			}
			conv_ts_filename := ConcatFilename(*basename, half_title, width.Widen.String(subtitle), ep_string, ".ts")
			converted_files = []string{conv_ts_filename, conv_ts_filename + ".err", conv_ts_filename + ".program.txt"}
			program_name = width.Widen.String(origin_program_name)
		}
	} else {
		converted_files = []string{fold_ts_filename, fold_ts_filename + ".err", fold_ts_filename + ".program.txt"}
		program_name = origin_program_name
	}

	err = OperateFile(save_path, origin_path, half_title, program_name, files, converted_files)
	if err != nil {
		return err
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
		cmd := exec.Command("C:\\Windows\\System32\\shutdown.exe", "/s", "/t", "60", "/f", "/c", "shutdown by epg-cycler after 60s")
		if cmd.Err != nil {
			return cmd.Err
		}
		DebugLog("execute shutdown")
		// Runだとシャットダウンし終えるまで処理が進まなくなるのでStartを使う
		if err := cmd.Start(); err != nil {
			return err
		}
		return nil
	}
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("shutdown", "-h now")
		if cmd.Err != nil {
			return cmd.Err
		}
		DebugLog("execute shutdown")
		if err := cmd.Start(); err != nil {
			return err
		}
		return nil
	}

	return UnSupoortedOSError()
}

func OperateFile(save_path, origin_path, title, program_name string, files []string, converted_names []string) error {
	program_save_path := filepath.Join(save_path, program_name)
	err := os.Mkdir(program_save_path, 0755)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	if strings.HasSuffix(program_save_path, " ") {
		program_save_path = program_save_path[:len(program_save_path)-1]
	}

	for idx, file := range files {
		err := os.Rename((filepath.Join(filepath.Dir(origin_path), file)), filepath.Join(program_save_path, converted_names[idx]))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Printf("%s: %s\n", os.ErrNotExist, filepath.Join(filepath.Dir(origin_path), converted_names[idx]))
				continue
			}
			return err
		}

	}
	return nil
}
