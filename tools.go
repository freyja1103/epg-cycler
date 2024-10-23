package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/width"
)

func GetProgramName(basename string) (string, string) {
	// 半角に変換して番組名を取得
	basename = width.Fold.String(basename)
	var end_brackets int = 0
	var exist bool
	var name string
	if strings.Index(basename, "[") == 0 {
		end_brackets = strings.Index(basename, "]") + 1
	}

	name, match, exist := GetNameByRegx(`第[0-9]+話`, basename)
	if !exist {
		name, match, exist = GetNameByRegx(`#[0-9]+`, basename)
	}
	if !exist {
		name, match, exist = GetNameByRegx(`第[0-9]+期`, basename)
	}
	if !exist {
		name, match, exist = GetNameByRegx(`([0-9]+)`, basename)
	}
	if !exist {
		name, match, exist = GetNameByRegx(`_`, basename)
	}
	if !exist {
		// たまにイレギュラーで ★最終話 みたいなのがあるので最終手段
		// タイトル内でスペース区切りの場合は対応してません
		log.Println(WarnProgramName())
		if strings.Index(basename, " ") == -1 {
			return basename[end_brackets:], match
		}
		return basename[end_brackets:strings.Index(basename, " ")], match
	}
	return name, match
}

func GetSubtitle(s string) (string, error) {
	regex := regexp.MustCompile(`[「」]`)
	matches := regex.FindAllStringIndex(s, -1)
	l := len(matches)
	if l > 0 {
		if l == 2 {
			return width.Widen.String(s[matches[0][1]:matches[1][0]]), nil
		}
		if l%4 == 0 {
			return width.Widen.String(s[matches[l-2][1]:matches[l-1][0]]), nil
		}
	}

	var match_ep []int

	match_date := regex_date.FindStringIndex(s)

	if regex_ep_kanji.MatchString(s) {
		match_ep = regex_ep_kanji.FindStringIndex(s)
		return strings.TrimSpace(s[match_ep[1]:match_date[0]]), nil
	}
	if regex_ep_shrp.MatchString(s) {
		match_ep = regex_ep_shrp.FindStringIndex(s)
		return strings.TrimSpace(s[match_ep[1]:match_date[0]]), nil
	}
	if regex_ep_brkt.MatchString(s) {
		match_ep = regex_ep_shrp.FindStringIndex(s)
		log.Println(match_ep, match_date)
		return strings.TrimSpace(s[match_ep[1]:match_date[0]]), nil
	}
	return "", SubtitleNotFoundError(s)
}

func GetNameByRegx(expr, basename string) (string, string, bool) {
	var name string
	if expr == `第[0-9]+期` {
		regex := regexp.MustCompile(expr)
		if regex.MatchString(basename) {
			matches := regex.FindAllString(basename, -1)
			name := basename[:strings.Index(basename, matches[0])+7]
			log.Println("reg: ", matches, name)
			return name, matches[0], true
		}
	}

	regex := regexp.MustCompile(expr)
	if regex.MatchString(basename) {
		matches := regex.FindAllString(basename, -1)
		name := basename[:strings.Index(basename, matches[0])]
		log.Println("reg: ", matches, name)
		return name, matches[0], true
	}
	return name, "", false
}

func isInvalidName(s string) (bool, []string) {
	s = width.Fold.String(s)
	invalidChars := `[<>:"/\|?*]`
	regex := regexp.MustCompile(invalidChars)
	regex.FindAllString(s, -1)
	if regex.MatchString(s) {
		return true, regex.FindAllString(s, -1)
	}
	return false, nil
}

func GetEpisodeName(s string) string {
	if regex_ep_kanji.MatchString(s) {
		return regex_ep_kanji.FindString(s)
	}
	if regex_ep_shrp.MatchString(s) {
		return regex_ep_shrp.FindString(s)
	}
	if regex_ep_brkt.MatchString(s) {
		return regex_ep_brkt.FindString(s)
	}
	return ""
}

func ConcatFilename(filename, title, subtitle, episode_string, ext string) string {
	subtitle_and_brackets := "「" + subtitle + "」"
	if subtitle == "" {
		subtitle_and_brackets = ""
	}
	return title + " " + episode_string + subtitle_and_brackets + regexp.MustCompile("_[0-9]{8}").FindString(width.Fold.String(filename)) + ext
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

func TidyAllFiles(save_path string) error {
	var (
		conv_filename     string
		filename          string
		program_name      string
		program_save_path string
	)
	files, err := SearchNotTidyFiles(save_path)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename = filepath.Base(file)
		folded_program_name, ep_string := GetProgramName(filename)

		subtitle, err := GetSubtitle(filename)
		if err != nil {
			if errors.Is(err, SubtitleNotFoundError(filename)) {
				subtitle = ""
			}
		}

		isInvalid, _ := isInvalidName(width.Fold.String(filename))
		isInvalidSubtitle, _ := isInvalidName(subtitle)
		isInvalidProgramName, _ := isInvalidName(folded_program_name)

		if isInvalid {
			DebugLog("Will be invalid filename, no convert fold style")
			if !isInvalidProgramName && isInvalidSubtitle {
				DebugLog("Only subtitle is invalid")
				// only subtitle is invalid

				conv_filename = ConcatFilename(filename, folded_program_name, width.Widen.String(subtitle), ep_string, filepath.Ext(file))
				program_name = folded_program_name
			} else {
				DebugLog("Program name is invalid")

				conv_filename = ConcatFilename(filename, width.Widen.String(folded_program_name), width.Fold.String(subtitle), ep_string, filepath.Ext(file))
				program_name = width.Widen.String(folded_program_name)
			}
		} else {
			conv_filename = width.Fold.String(filename)
			program_name = folded_program_name
		}

		program_save_path = filepath.Join(save_path, program_name)
		if strings.HasSuffix(program_save_path, " ") {
			program_save_path = program_save_path[:len(program_save_path)-1]
		}
		log.Println("save path: ", program_save_path)
		err = os.Mkdir(program_save_path, 0755)
		if err != nil && !os.IsExist(err) {
			Errorlog(err)
		}

		save := filepath.Join(program_save_path, conv_filename)
		log.Println("move to :", save)
		err = os.Rename(file, save)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Printf("The file or directory does not exist: %s, %s\n", file, save)
				continue
			}
			Errorlog(err)
			continue
		}
	}
	return nil
}

type targetProcesses []string

func (tp *targetProcesses) String() string {
	return fmt.Sprintf("%v", *tp)
}

func (tp *targetProcesses) Set(value string) error {
	*tp = append(*tp, value)
	return nil
}
