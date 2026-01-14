package edcbapiclient

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/freyja1103/epg-cycler/edcb-api-client/dto"
)

type API interface {
	GetEnumReserveInfo(ctx context.Context) (*dto.ReserveInfoEntry, error)
	GetEnumRecInfo(ctx context.Context, params *EnumRecInfoParams) ([]*RecInfo, error)
}

func NewEDCBAPIClient(hostname string, client *http.Client) API {
	if client == nil {
		client = http.DefaultClient
	}
	return &api{hostname: hostname, client: client}
}

type api struct {
	hostname string
	client   *http.Client
}

func (a *api) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}

func (a *api) BaseURL() string {
	return fmt.Sprintf("http://%s/api", a.hostname)
}

func (a *api) GetEnumReserveInfo(ctx context.Context) (*dto.ReserveInfoEntry, error) {
	dest := fmt.Sprintf("%s/EnumReserveInfo", a.BaseURL())
	res, err := a.get(ctx, dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get EnumReserveInfo: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	entry := new(dto.ReserveInfoEntry)
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml: %w", err)
	}

	return entry, nil
}

type EnumRecInfoParams struct {
	ID string
}

func (a *api) GetEnumRecInfo(ctx context.Context, params *EnumRecInfoParams) ([]*RecInfo, error) {
	dest := fmt.Sprintf("%s/EnumRecInfo", a.BaseURL())
	uv := url.Values{}
	if params != nil {
		uv.Add("id", params.ID)
	}
	res, err := a.get(ctx, dest+"?"+uv.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to get EnumRecInfo: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	entry := new(dto.RecInfoEntry)
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml: %w", err)
	}

	return toRecInfo(entry), nil
}

func toRecInfo(entry *dto.RecInfoEntry) []*RecInfo {
	recInfos := make([]*RecInfo, 0, len(entry.Items.RecInfos))
	for _, r := range entry.Items.RecInfos {
		startTime, err := parseDateTime(r.StartDate, r.StartTime)
		if err != nil {
			slog.Error("failed to parse date time", slog.Any("error", err))
			continue
		}
		recInfos = append(recInfos, &RecInfo{
			ID:          r.ID,
			ServiceID:   r.SID,
			ServiceName: r.ServiceName,
			Duration:    r.Duration,
			StartTime:   startTime,
			RecFilePath: r.RecFilePath,
			Title:       r.Title,
		})
	}
	return recInfos
}

func parseDateTime(dateStr, timeStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	return time.ParseInLocation(
		"2006/01/02 15:04:05",
		dateStr+" "+timeStr,
		loc,
	)
}
