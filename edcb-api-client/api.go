package edcbapiclient

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type API interface {
	GetEnumReserveInfo(ctx context.Context) (*ReserveInfoEntry, error)
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

func (a *api) GetEnumReserveInfo(ctx context.Context) (*ReserveInfoEntry, error) {
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

	entry := new(ReserveInfoEntry)
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml: %w", err)
	}

	return entry, nil
}
