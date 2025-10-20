package epgcycler

import (
	"strings"

	"golang.org/x/text/width"
)

func getSubtitle(s string) (string, error) {
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

func convertInvalidCharsToWiden(s string) string {
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
