package validation

import (
	"regexp"

	"golang.org/x/text/width"
)

var invalidCharsRegex = regexp.MustCompile(`[<>:"/\|?*]`)

func IsInvalidName(s string) (bool, []string) {
	s = width.Fold.String(s)
	invalidCharsRegex.FindAllString(s, -1)
	if invalidCharsRegex.MatchString(s) {
		return true, invalidCharsRegex.FindAllString(s, -1)
	}
	return false, nil
}

func IsValidName(s string) (bool, []string) {
	s = width.Fold.String(s)
	invalidCharsRegex.FindAllString(s, -1)
	if !invalidCharsRegex.MatchString(s) {
		return false, nil
	}
	return true, invalidCharsRegex.FindAllString(s, -1)
}
