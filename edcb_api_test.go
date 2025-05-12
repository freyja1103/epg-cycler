package main

import (
	"context"
	"encoding/xml"
	"log/slog"
	"testing"

	"github.com/freyja1103/epg-cycler/logging"
)

func TestEDCBAPI(t *testing.T) {
	body, err := GetEnumReserveInfo("localhost:5510")
	if err != nil {
		logging.Error(err.Error())
	}

	var entry Entry
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		logging.Error(err.Error())
	}

	hasReserve, timeList, err := HasRemainReserve(&entry)
	if hasReserve {
		if err != nil {
			logging.Error(err.Error())
		}
		logging.InfoAttrs(context.Background(), slog.LevelInfo, "test EDCB API (EnumReserveInfo)", slog.Any("timeList", timeList), slog.Bool("hasReserve", hasReserve))

	}
}
