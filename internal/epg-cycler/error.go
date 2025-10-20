package epgcycler

import "errors"

var (
	ErrSubtitleNotFound   = errors.New("subtitle not found")
	ErrInvalidProgramFile = errors.New("invalid .ts.program.txt format")
)
