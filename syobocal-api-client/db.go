package syobocalapiclient

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"github.com/freyja1103/epg-cycler/syobocal-api-client/dto"
)

type ProgLookupParams struct {
	TIDs       []string
	ChIDs      []string
	Range      *Range
	Counts     []string
	LastUpdate *Range
	Fields     []Field
	PIDs       []string
}

func (a *api) ProgLookup(ctx context.Context, params *ProgLookupParams) ([]*ProgItem, error) {
	if params == nil {
		return nil, fmt.Errorf("params must be set")
	}

	uv := url.Values{}
	uv.Add("Command", "ProgLookup")
	uv.Add("TID", toCommaSeparatedString(params.TIDs))
	uv.Add("ChID", toCommaSeparatedString(params.ChIDs))

	if params.Range != nil {
		if params.Range.To == nil {
			uv.Add("StTime", params.Range.toParam())
		} else {
			uv.Add("Range", params.Range.toParam())
		}
	}

	uv.Add("Count", toCommaSeparatedString(params.Counts))

	if params.LastUpdate != nil {
		if params.LastUpdate.To != nil {
			return nil, fmt.Errorf("LastUpdate must have To set")
		}

		uv.Add("LastUpdate", params.LastUpdate.toParam())
	}

	uv.Add("Fields", fieldsToParam(params.Fields))
	uv.Add("PID", toCommaSeparatedString(params.PIDs))

	res, err := a.get(ctx, fmt.Sprintf("%s/%s", a.baseURL, endopointDb), uv)
	if err != nil {
		return nil, fmt.Errorf("failed to get ProgLookup: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	r := new(dto.ProgLookupResponse)
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml: %w", err)
	}

	if r.Result.Code != http.StatusOK {
		return nil, fmt.Errorf("API error: %d %s", r.Result.Code, r.Result.Message)
	}

	return ToProgItem(r.ProgItems.ProgItem), nil
}

func ToProgItem(items []*dto.ProgItem) []*ProgItem {
	p := make([]*ProgItem, 0, len(items))
	for _, item := range items {
		p = append(p, &ProgItem{
			ID:          item.ID,
			PID:         item.PID,
			TID:         item.TID,
			StOffset:    item.StOffset,
			Count:       item.Count,
			Flag:        item.Flag,
			Deleted:     item.Deleted,
			Warn:        item.Warn,
			ChID:        item.ChID,
			Revision:    item.Revision,
			StTime:      item.StTime,
			SubTitle:    item.SubTitle,
			ProgComment: item.ProgComment,
			EdTime:      item.EdTime,
			LastUpdate:  item.LastUpdate,
		})
	}
	return p
}

type TitleLookupParams struct {
	TIDs       []string
	LastUpdate *Range
	Fields     []Field
}

func (a *api) TitleLookup(ctx context.Context, params *TitleLookupParams) ([]*TitleItem, error) {
	if params == nil {
		return nil, fmt.Errorf("params must be set")
	}

	if slices.Contains(params.TIDs, "*") {
		return nil, fmt.Errorf("wildcard '*' is not supported for this api client")
	}

	uv := url.Values{}
	uv.Add("Command", "TitleLookup")
	uv.Add("TID", toCommaSeparatedString(params.TIDs))
	res, err := a.get(context.Background(), fmt.Sprintf("%s/%s", a.baseURL, endopointDb), uv)
	if err != nil {
		return nil, fmt.Errorf("failed to get TitleLookup: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	r := new(dto.TitleLookupResponse)
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml: %w", err)
	}

	if r.Result.Code != http.StatusOK {
		return nil, fmt.Errorf("API error: %d %s", r.Result.Code, r.Result.Message)
	}
	return toTitleItem(r.TitleItems.TitleItem), nil
}

func toTitleItem(items []*dto.TitleItem) []*TitleItem {
	t := make([]*TitleItem, 0, len(items))
	for _, item := range items {
		t = append(t, &TitleItem{
			ID:            item.ID,
			TID:           item.TID,
			LastUpdate:    item.LastUpdate,
			Title:         item.Title,
			ShortTitle:    item.ShortTitle,
			TitleYomi:     item.TitleYomi,
			TitleEN:       item.TitleEN,
			Comment:       item.Comment,
			Cat:           item.Cat,
			TitleFlag:     item.TitleFlag,
			FirstYear:     item.FirstYear,
			FirstMonth:    item.FirstMonth,
			FirstEndYear:  item.FirstEndYear,
			FirstEndMonth: item.FirstEndMonth,
			FirstCh:       item.FirstCh,
			Keywords:      item.Keywords,
			UserPoint:     item.UserPoint,
			UserPointRank: item.UserPointRank,
			SubTitles:     item.SubTitles,
		})
	}
	return t
}
