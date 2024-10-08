package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/text/width"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
	procPreventShutdown := flag.String(args[5], "", "process that prevent shutdown")
	APIURL := flag.String(args[6], "localhost:5510", "EpgTimer's HTTP server, IP:port")
	all_tidy_mode := flag.Bool(args[7], false, "The mode for sorting all the recording files in the directory at once")
	flag.Parse()

	save_path := strings.ReplaceAll(filepath.ToSlash(*s_path), "//", "/")
	origin_path := strings.ReplaceAll(filepath.ToSlash(*o_path), "//", "/")

	// 一括整理モード isInvalidName使ってないので多分まだ動かない
	if *all_tidy_mode {
		files, err := SearchNotTidyFiles(save_path)
		if err != nil {
			Errorlog(err)
			return
		}
		for _, file := range files {
			program_name := GetProgramName(filepath.Base(file))
			program_save_path := filepath.Join(save_path, program_name)

			if strings.HasSuffix(program_save_path, " ") {
				program_save_path = program_save_path[:len(program_save_path)-1]
			}
			err := os.Mkdir(program_save_path, 0755)
			if err != nil && !os.IsExist(err) {
				Errorlog(err)
			}

			save := width.Fold.String(filepath.Join(program_save_path, filepath.Base(file)))
			log.Println("From: ", file)
			log.Println("To: ", save)
			err = os.Rename(file, save)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					log.Printf("The file or directory does not exist: %s, %s\n", file, save)
					continue
				}
				Errorlog(err)
				log.Println("May be invaild filename, skipped: ", save)
				continue
			}
		}
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

}

func SearchNotTidyFiles(save_path string) ([]string, error) {
	files := []string{}
	filepath.WalkDir(save_path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		path = filepath.ToSlash(path)

		if filepath.ToSlash(filepath.Dir(path)) == save_path && !info.Type().IsDir() {
			if strings.Contains(filepath.Ext(path), ".ts") || strings.Contains(filepath.Ext(path), ".err") || strings.Contains(filepath.Ext(path), ".txt") {
				files = append(files, path)
			}

		}
		return nil
	})
	for _, v := range files {
		log.Println(v)
	}
	return files, nil
}

func dirTidy(save_path, origin_path string, title, basename *string) error {
	// basename -> $FileName$
	var ts_filename string = *basename + ".ts"
	var half_title string
	var subtitle string
	var err error
	var files []string
	var fold_ts_filename = width.Fold.String(*basename + ".ts")
	half_title = width.Fold.String(*title)

	isInvalid, _ := isInvalidName(*basename)
	if isInvalid {
		DebugLog("Will be invalid filename, no convert fold style")
		subtitle, err = GetSubtitle(fold_ts_filename)
		if err != nil {
			return err
		}

		origin_program_name := GetProgramName(*basename)
		isInvalidSubtitle, _ := isInvalidName(subtitle)
		isInvalidProgramName, _ := isInvalidName(origin_program_name)
		if !isInvalidProgramName && isInvalidSubtitle {
			// only subtitle is invalid
			program_name := GetProgramName(width.Fold.String(*basename))
			files = []string{ts_filename, ts_filename + ".err", ts_filename + ".program.txt"}
			conv_ts_filename := half_title + " " + regexp.MustCompile("#\\d+").FindString(fold_ts_filename) + "「" + subtitle + "」" + regexp.MustCompile("_[0-9]{8}").FindString(*basename) + ".ts"
			conv_files := []string{conv_ts_filename, conv_ts_filename + ".err", conv_ts_filename + ".program.txt"}
			err := OperateFile(save_path, origin_path, half_title, program_name, files, conv_files)
			if err != nil {
				return err
			}
			return nil
		}

		files = []string{ts_filename, ts_filename + ".err", ts_filename + ".program.txt"}
		if strings.HasSuffix(origin_program_name, " ") {
			origin_program_name = origin_program_name[:len(origin_program_name)-1]
		}
		widen_program_name := width.Widen.String(origin_program_name)
		err := OperateFile(save_path, origin_path, *title, widen_program_name, files, files)
		if err != nil {
			return err
		}
		return nil

	} else {
		files = []string{ts_filename, ts_filename + ".err", ts_filename + ".program.txt"}
		folded_files := []string{fold_ts_filename, fold_ts_filename + ".err", fold_ts_filename + ".program.txt"}
		program_name := GetProgramName(width.Fold.String(*basename))

		err := OperateFile(save_path, origin_path, half_title, program_name, files, folded_files)
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
		DebugLog("execute shutdown")
		return nil
	}
	if runtime.GOOS == "darwin" {
		exec := exec.Command("shutdown", "-h now")
		if exec.Err != nil {
			return exec.Err
		}
		DebugLog("execute shutdown")
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
