package syobocalapiclient

type ProgItem struct {
	ID          int
	PID         int
	TID         int
	StOffset    int
	Count       *int
	Flag        int
	Deleted     int
	Warn        int
	ChID        int
	Revision    int
	StTime      string
	SubTitle    string
	ProgComment string
	EdTime      string
	LastUpdate  string
}
