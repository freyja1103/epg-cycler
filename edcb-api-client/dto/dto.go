package dto

import "encoding/xml"

type ReserveInfoEntry struct {
	XMLName xml.Name         `xml:"entry"`
	Total   int              `xml:"total"`
	Index   int              `xml:"index"`
	Count   int              `xml:"count"`
	Items   ReserveInfoItems `xml:"items"`
}

type ReserveInfoItems struct {
	ReserveInfo []*ReserveInfo `xml:"reserveinfo"`
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

type RecInfoEntry struct {
	XMLName xml.Name     `xml:"entry"`
	Total   int          `xml:"total"`
	Index   int          `xml:"index"`
	Count   int          `xml:"count"`
	Items   RecInfoItems `xml:"items"`
}

type RecInfoItems struct {
	RecInfos []*RecInfo `xml:"recinfo"`
}

type RecInfo struct {
	ID                int    `xml:"ID"`
	RecFilePath       string `xml:"recFilePath"`
	Title             string `xml:"title"`
	ONID              int    `xml:"ONID"`
	TSID              int    `xml:"TSID"`
	SID               int    `xml:"SID"`
	EventID           int    `xml:"eventID"`
	ServiceName       string `xml:"service_name"`
	StartDate         string `xml:"startDate"`
	StartTime         string `xml:"startTime"`
	StartDayOfWeek    int    `xml:"startDayOfWeek"`
	Duration          int    `xml:"duration"`
	StartDateEpg      string `xml:"startDateEpg"`
	StartTimeEpg      string `xml:"startTimeEpg"`
	StartDayOfWeekEpg int    `xml:"startDayOfWeekEpg"`
	Drops             int    `xml:"drops"`
	Scrambles         int    `xml:"scrambles"`
	RecStatus         int    `xml:"recStatus"`
	Comment           string `xml:"comment"`
	ProgramInfo       string `xml:"programInfo"`
	ErrInfo           string `xml:"errInfo"`
	Protect           int    `xml:"protect"`
}
