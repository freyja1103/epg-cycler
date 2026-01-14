package syobocalapiclient

type TitleItem struct {
	ID            int
	TID           int
	LastUpdate    string
	Title         string
	ShortTitle    string
	TitleYomi     string
	TitleEN       string
	Comment       string
	Cat           int
	TitleFlag     int
	FirstYear     int
	FirstMonth    int
	FirstEndYear  *int
	FirstEndMonth *int
	FirstCh       string
	Keywords      string
	UserPoint     int
	UserPointRank int
	SubTitles     string
}
