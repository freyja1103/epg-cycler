package model

import "encoding/xml"

type TimeInfo struct {
	StartDate string
	StartTime string
}

type ReserveProgram struct {
	ID             int
	Title          string
	StartDate      string
	StartTime      string
	StartDayOfWeek int
	Duration       int
}

type ReservedProgramRepository interface {
	Get(url string) *Entry
	HasRemainReserve(rsrv ReserveInfo) bool
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Total   int      `xml:"total"`
	Index   int      `xml:"index"`
	Count   int      `xml:"count"`
	Items   Items    `xml:"items"`
}

type Items struct {
	ReserveInfo []ReserveInfo `xml:"reserveinfo"`
}

type ReserveInfo struct {
	ID              int        `xml:"ID"`
	Title           string     `xml:"title"`
	StartDate       string     `xml:"startDate"`
	StartTime       string     `xml:"startTime"`
	StartDayOfWeek  int        `xml:"startDayOfWeek"`
	Duration        int        `xml:"duration"`
	ServiceName     string     `xml:"service_name"`
	ONID            int        `xml:"ONID"`
	TSID            int        `xml:"TSID"`
	SID             int        `xml:"SID"`
	EventID         int        `xml:"eventID"`
	Comment         string     `xml:"comment"`
	OverlapMode     int        `xml:"overlapMode"`
	RecSetting      RecSetting `xml:"recsetting"`
	RecFileNameList []string   `xml:"recFileNameList>recFileName"`
}

type RecSetting struct {
	RecMode          int      `xml:"recMode"`
	Priority         int      `xml:"priority"`
	TuijyuuFlag      int      `xml:"tuijyuuFlag"`
	ServiceMode      int      `xml:"serviceMode"`
	BatFilePath      string   `xml:"batFilePath"`
	RecFolderList    []string `xml:"recFolderList>folder"`
	SuspendMode      int      `xml:"suspendMode"`
	DefserviceMode   int      `xml:"defserviceMode"`
	RebootFlag       int      `xml:"rebootFlag"`
	UseMargineFlag   int      `xml:"useMargineFlag"`
	StartMargine     int      `xml:"startMargine"`
	EndMargine       int      `xml:"endMargine"`
	ContinueRecFlag  int      `xml:"continueRecFlag"`
	PartialRecFlag   int      `xml:"partialRecFlag"`
	TunerID          int      `xml:"tunerID"`
	PartialRecFolder []string `xml:"partialRecFolder>folder"`
}
