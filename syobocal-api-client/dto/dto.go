package dto

import "encoding/xml"

type TitleLookupResponse struct {
	XMLName    xml.Name          `xml:"TitleLookupResponse"`
	Result     TitleLookupResult `xml:"Result"`
	TitleItems TitleItems        `xml:"TitleItems"`
}

type TitleLookupResult struct {
	Code    int    `xml:"Code"`
	Message string `xml:"Message"`
}

type TitleItems struct {
	TitleItem []*TitleItem `xml:"TitleItem"`
}

type TitleItem struct {
	ID            int    `xml:"id,attr"`
	TID           int    `xml:"TID"`
	LastUpdate    string `xml:"LastUpdate"`
	Title         string `xml:"Title"`
	ShortTitle    string `xml:"ShortTitle"`
	TitleYomi     string `xml:"TitleYomi"`
	TitleEN       string `xml:"TitleEN"`
	Comment       string `xml:"Comment"`
	Cat           int    `xml:"Cat"`
	TitleFlag     int    `xml:"TitleFlag"`
	FirstYear     int    `xml:"FirstYear"`
	FirstMonth    int    `xml:"FirstMonth"`
	FirstEndYear  *int   `xml:"FirstEndYear"`  // 空要素あり
	FirstEndMonth *int   `xml:"FirstEndMonth"` // 空要素あり
	FirstCh       string `xml:"FirstCh"`
	Keywords      string `xml:"Keywords"`
	UserPoint     int    `xml:"UserPoint"`
	UserPointRank int    `xml:"UserPointRank"`
	SubTitles     string `xml:"SubTitles"`
}

type ProgLookupResponse struct {
	XMLName   xml.Name         `xml:"ProgLookupResponse"`
	ProgItems ProgItems        `xml:"ProgItems"`
	Result    ProgLookupResult `xml:"Result"`
}

type ProgItems struct {
	ProgItem []*ProgItem `xml:"ProgItem"`
}

type ProgItem struct {
	ID          int    `xml:"id,attr"`
	LastUpdate  string `xml:"LastUpdate"`
	PID         int    `xml:"PID"`
	TID         int    `xml:"TID"`
	StTime      string `xml:"StTime"`
	StOffset    int    `xml:"StOffset"`
	EdTime      string `xml:"EdTime"`
	Count       *int   `xml:"Count"` // 空要素あり
	SubTitle    string `xml:"SubTitle"`
	ProgComment string `xml:"ProgComment"`
	Flag        int    `xml:"Flag"`
	Deleted     int    `xml:"Deleted"`
	Warn        int    `xml:"Warn"`
	ChID        int    `xml:"ChID"`
	Revision    int    `xml:"Revision"`
}

type ProgLookupResult struct {
	Code    int    `xml:"Code"`
	Message string `xml:"Message"`
}
