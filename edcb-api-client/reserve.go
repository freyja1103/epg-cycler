package edcbapiclient

type ReserveInfo struct {
	ID             int
	Title          string
	StartDate      string
	StartTime      string
	StartDayOfWeek int
	Duration       int
	ServiceName    string
	ONID           int
	TSID           int
	SID            int
	EventID        int
}
