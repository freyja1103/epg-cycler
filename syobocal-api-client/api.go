package syobocalapiclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type api struct {
	client  *http.Client
	baseURL string
}

func NewSyobocalAPIClient(client *http.Client) API {
	if client == nil {
		client = http.DefaultClient
	}
	return &api{client: client, baseURL: "https://cal.syoboi.jp"}
}

type API interface {
	ProgLookup(ctx context.Context, params *ProgLookupParams) ([]*ProgItem, error)
	TitleLookup(ctx context.Context, params *TitleLookupParams) ([]*TitleItem, error)
}

func (a *api) get(ctx context.Context, url string, params url.Values) (*http.Response, error) {
	return a.call(ctx, http.MethodGet, url+"?"+params.Encode())
}

func (a *api) call(ctx context.Context, method string, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	// req.Header.Set("User-Agent", "syobocal-api-client (+https://github.com/freyja1103/epg-cycler)")
	res, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %v", err)
	}
	defer res.Body.Close()

	return res, nil
}

const (
	endopointDb = "db.php"
)

type Field string

func fieldsToParam(fields []Field) string {
	strFields := make([]string, 0, len(fields))
	for _, f := range fields {
		strFields = append(strFields, string(f))
	}
	return strings.Join(strFields, ",")
}

const (
	FieldTID          Field = "TID"
	FieldTitle        Field = "Title"
	FieldTitleYomi    Field = "TitleYomi"
	FieldComment      Field = "Comment"
	FieldCat          Field = "Cat"
	FieldTitleFlag    Field = "TitleFlag"
	FieldFirstYear    Field = "FirstYear"
	FieldFirstMonth   Field = "FirstMonth"
	FieldFirstEndYear Field = "FirstEndYear"
	FirstEndMonth     Field = "FirstEndMonth"
	FirstCh           Field = "FirstCh"
	FieldSubTitles    Field = "SubTitles"
)

type Range struct {
	From time.Time
	To   *time.Time
}

const timeLayout = "20060102_150405"

func (r *Range) toParam() string {
	if r.To == nil {
		return fmt.Sprintf("%s-", r.From.Format(timeLayout))
	}
	return fmt.Sprintf("%s-%s", r.From.Format(timeLayout), r.To.Format(timeLayout))
}

func toCommaSeparatedString(arr []string) string {
	return strings.Join(arr, ",")
}
