package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

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

type TimeInfo struct {
	StartDate string
	StartTime string
}

func APIReq2Body(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HasRemainReserve(entry *Entry) (bool, error) {
	now := time.Now()
	date := (now.Format("2006/01/02"))
	var timeList []TimeInfo
	for _, reserve := range entry.Items.ReserveInfo {
		startDate, startTime, err := ParseTime(reserve.StartDate, reserve.StartTime)
		if err != nil {
			return true, err
		}
		// 同じ日付または次の日の4時までの予約を対象とする
		if 0 <= now.Hour() || now.Hour() <= 4 {
			if startDate.Day() == now.Day() && startTime.Hour() < 4 {
				timeList = append(timeList, TimeInfo{
					reserve.StartDate,
					reserve.StartTime,
				})
			}
		}
		if now.Hour() > 4 {
			if date == reserve.StartDate {
				timeList = append(timeList, TimeInfo{
					reserve.StartDate,
					reserve.StartTime,
				})
			}
		}
	}

	return len(timeList) >= 1, nil
}

func ParseTime(d string, t string) (time.Time, time.Time, error) {
	zeroTime := time.Time{}
	startDate, err := time.Parse("2006/01/02", d)
	if err != nil {
		return zeroTime, zeroTime, err
	}
	startTime, err := time.Parse("15:04:05", t)
	if err != nil {
		return zeroTime, zeroTime, err
	}
	return startDate, startTime, nil
}
