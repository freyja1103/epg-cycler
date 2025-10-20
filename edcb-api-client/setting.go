package edcbapiclient

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
