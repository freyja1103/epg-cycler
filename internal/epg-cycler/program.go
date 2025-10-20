package epgcycler

import (
	"context"
	"log/slog"
	"regexp"
	"strings"

	"golang.org/x/text/width"
)

type Program struct {
	Title       string
	Episode     string
	Subtitle    string
	ProgramName string
}

func NewProgram(programName string) *Program {
	return &Program{
		Title:       "",
		Subtitle:    "",
		Episode:     "",
		ProgramName: programName,
	}
}

func SplitProgramName(fullname string) (title string, episodeAndSubtitle string) {
	parts := strings.Split(fullname, "　")
	if len(parts) < 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

func (p *Program) SetTitle(title string) {
	p.Title = title
}

func (p *Program) SetEpisode(episode string) {
	p.Episode = width.Fold.String(episode)
}

func (p *Program) SetSubtitle(subtitle string) {
	p.Subtitle = width.Fold.String(subtitle)
}

func (p *Program) GetTitle(ctx context.Context) (title string, episode string) {
	programName := p.ProgramName
	var end_brackets int = 0
	// Regex patternじゃなくて.program.txtからparseして取得するほうがよさそう
	name, match, exist := p.GetNameByRegx(ctx, episodeRegex)
	if !exist {
		name, match, exist = p.GetNameByRegx(ctx, underlineRegex)
	}
	if !exist {
		// たまにイレギュラーで ★最終話 みたいなのがあるので最終手段
		// タイトル内でスペース区切りの場合は対応してません
		slog.InfoContext(ctx, "the format of the program name is not supported. the name of the created directory may differ from the actual program name.")
		if strings.Index(programName, "[") == 0 {
			end_brackets = strings.Index(programName, "]") + 1
		}
		if !strings.Contains(programName, " ") {
			return programName[end_brackets:], match
		}
		return programName[end_brackets:strings.Index(programName, " ")], match
	}
	name, _ = strings.CutSuffix(name, " ")

	p.SetTitle(name)

	return name, match
}

func (p *Program) GetNameByRegx(ctx context.Context, regex *regexp.Regexp) (name string, matchedName string, _ bool) {
	if regex.MatchString(p.ProgramName) {
		matches := regex.FindAllString(p.ProgramName, -1)
		name := p.ProgramName[:strings.Index(p.ProgramName, matches[0])]
		slog.InfoContext(ctx, "reg", slog.Any("matches", matches), slog.String("name", name))
		return name, matches[0], true
	}
	return "", "", false
}
