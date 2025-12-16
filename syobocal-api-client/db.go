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

func (a *api) ProgLookup(ctx context.Context, params *ProgLookupParams) ([]*dto.ProgItem, error) {
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

	return r.ProgItems.ProgItem, nil
}

type TitleLookupParams struct {
	TIDs       []string
	LastUpdate *Range
	Fields     []Field
}

func (a *api) TitleLookup(params *TitleLookupParams) ([]*dto.TitleItem, error) {
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
	return r.TitleItems.TitleItem, nil
}
