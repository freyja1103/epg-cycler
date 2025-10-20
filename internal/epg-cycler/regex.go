package epgcycler

import "regexp"

var (
	regexp_episode = regexp.MustCompile(`_[0-9]{8}|第[0-9]+話|#[0-9]+|#[0-9]+|([0-9]{,2})`)
	regexp_date    = regexp.MustCompile(`_[0-9]{8}`)

	episodeRegex      = regexp.MustCompile(`_[0-9]{8}|第[0-9]+話|#[0-9]+|#[0-9]+|([0-9]{,2})`)
	dateRegex         = regexp.MustCompile(`_[0-9]{8}`)
	underlineRegex    = regexp.MustCompile(`_`)
	quoteRegex        = regexp.MustCompile(`[「」]`)
	invalidCharsRegex = regexp.MustCompile(`[<>:"/\|?*]`)
)
