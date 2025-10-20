package edcbapiclient

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
