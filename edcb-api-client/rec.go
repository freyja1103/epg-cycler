package edcbapiclient

import "time"

type RecInfo struct {
	ID          int
	ServiceID   int
	Duration    int
	StartTime   time.Time
	RecFilePath string
	Title       string
	ServiceName string
}
