package main

import (
	"log"
	"regexp"
	"strings"

	"golang.org/x/text/width"
)

func GetProgramName(basename string) string {
	// 半角に変換して番組名を取得
	basename = width.Fold.String(basename)
	var end_brackets int = 0
	var exist bool
	var name string
	if strings.Index(basename, "[") == 0 {
		end_brackets = strings.Index(basename, "]") + 1
	}

	name, exist = GetNameByRegx(`第[0-9]+話`, basename)
	if !exist {
		name, exist = GetNameByRegx(`#[0-9]+`, basename)
	}
	if !exist {
		name, exist = GetNameByRegx(`第[0-9]+期`, basename)
	}
	if !exist {
		name, exist = GetNameByRegx(`_`, basename)
	}
	if !exist {
		// たまにイレギュラーで ★最終話 みたいなのがあるので最終手段
		// タイトル内でスペース区切りの場合は対応してません
		log.Println(WarnProgramName())
		if strings.Index(basename, " ") == -1 {
			return basename[end_brackets:]
		}
		return basename[end_brackets:strings.Index(basename, " ")]
	}
	return name
}

func GetSubtitle(s string) (string, error) {
	regex := regexp.MustCompile(`[「」]`)
	matches := regex.FindAllStringIndex(s, -1)
	l := len(matches)
	if l == 2 {
		return width.Widen.String(s[matches[0][1]:matches[1][0]]), nil
	}
	if l%4 == 0 {
		return width.Widen.String(s[matches[l-2][1]:matches[l-1][0]]), nil
	}

	return "", SubtitleNotFoundError(s)
}

func GetNameByRegx(expr, basename string) (string, bool) {
	var name string
	if expr == `第[0-9]+期` {
		regex := regexp.MustCompile(expr)
		if regex.MatchString(basename) {
			matches := regex.FindAllString(basename, -1)
			name := basename[:strings.Index(basename, matches[0])+7]
			log.Println("reg: ", matches, name)
			return name, true
		}
	}

	regex := regexp.MustCompile(expr)
	if regex.MatchString(basename) {
		matches := regex.FindAllString(basename, -1)
		name := basename[:strings.Index(basename, matches[0])]
		log.Println("reg: ", matches, name)
		return name, true
	}
	return name, false
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
