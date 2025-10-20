package main

import (
	"bufio"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/freyja1103/epg-cycler/logging"
	"github.com/freyja1103/epg-cycler/program"
	"golang.org/x/text/width"
)

type ProgramOperation interface {
}

type programOperation struct {
	program program.Program
}

func NewProgramOperation(program program.Program) ProgramOperation {
	return &programOperation{program: program}
}

func (p *programOperation) GetTitle() (string, string) {
	basename := p.program.GetBasename()
	var end_brackets int = 0
	// Regex patternじゃなくて.program.txtからparseして取得するほうがよさそう
	name, match, exist := p.GetNameByRegx(episodeRegex)
	if !exist {
		name, match, exist = p.GetNameByRegx(underlineRegex)
	}
	if !exist {
		// たまにイレギュラーで ★最終話 みたいなのがあるので最終手段
		// タイトル内でスペース区切りの場合は対応してません
		logging.Warn("the format of the program name is not supported. the name of the created directory may differ from the actual program name.")
		if strings.Index(basename, "[") == 0 {
			end_brackets = strings.Index(basename, "]") + 1
		}
		if !strings.Contains(basename, " ") {
			return basename[end_brackets:], match
		}
		return basename[end_brackets:strings.Index(basename, " ")], match
	}
	name, _ = strings.CutSuffix(name, " ")

	p.program.SetTitle(name)

	return name, match
}

func GetProgramNameFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var thirdLine string
	for scanner.Scan() {
		lineNum++
		if lineNum == 3 {
			thirdLine = scanner.Text()
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if thirdLine == "" {
		return "", nil
	}
	parts := strings.Split(thirdLine, "　") // 全角スペース
	if len(parts) > 1 {
		return parts[1], nil
	}
	return thirdLine, nil
}

func (p *programOperation) GetNameByRegx(regex *regexp.Regexp) (string, string, bool) {
	if regex.MatchString(p.program.GetBasename()) {
		matches := regex.FindAllString(p.program.GetBasename(), -1)
		name := p.program.GetBasename()[:strings.Index(p.program.GetBasename(), matches[0])]
		logging.Info("reg", slog.Any("matches", matches), slog.String("name", name))
		return name, matches[0], true
	}
	return "", "", false
}

func NewProgram(basename, title string) (*program.Program, error) {
	subtitle, err := _GetSubtitle(basename)
	if err != nil {
		logging.Error("failed to get subtitle:", slog.Any("error", err))
		return nil, err
	}

	return &program.Program{
		Title:    ConvertInvalidCharsToWiden(width.Fold.String(title)),
		Subtitle: ConvertInvalidCharsToWiden(subtitle),
		Episode:  "", // Episode can be set later
		Basename: basename,
	}, nil
}

func _GetSubtitle(s string) (string, error) {
	matches := quoteRegex.FindAllStringIndex(s, -1)
	l := len(matches)
	if l == 2 {
		return width.Widen.String(s[matches[0][1]:matches[1][0]]), nil
	}
	if l%4 == 0 {
		return width.Widen.String(s[matches[l-2][1]:matches[l-1][0]]), nil
	}

	var match_ep []int
	date_idx := dateRegex.FindStringIndex(s)
	if episodeRegex.MatchString(s) {
		match_ep = episodeRegex.FindStringIndex(s)
		if len(match_ep) == 0 {
			return "", ErrSubtitleNotFound
		}
		return strings.TrimSpace(s[match_ep[0]:date_idx[1]]), nil
	}
	return "", ErrSubtitleNotFound
}

func ConvertInvalidCharsToWiden(s string) string {
	zenkakuMap := map[string]string{
		"<":  "＜",
		">":  "＞",
		":":  "：",
		`"`:  "＂",
		"/":  "／",
		"\\": "＼",
		"|":  "｜",
		"?":  "？",
		"*":  "＊",
	}
	return invalidCharsRegex.ReplaceAllStringFunc(s, func(m string) string {
		if z, ok := zenkakuMap[m]; ok {
			return z
		}
		return m
	})
}

func OperateDirectory(savePath, originPath, title, basename string) {
	program := &program.Program{
		Title: title,
	}
}
