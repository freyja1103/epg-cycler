package main

import "regexp"

var (
	regex_date     = regexp.MustCompile(`_[0-9]{8}`)
	regex_ep_kanji = regexp.MustCompile(`第[0-9]+話`)
	regex_ep_shrp  = regexp.MustCompile(`#[0-9]+`)
	regex_ep_brkt  = regexp.MustCompile(`([0-9]{,2})`)
)
