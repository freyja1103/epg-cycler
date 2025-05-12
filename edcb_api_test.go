package main

import (
	"context"
	"encoding/xml"
	"log/slog"
	"testing"

	"github.com/freyja1103/epg-cycler/logging"
)

func TestEDCBAPI(t *testing.T) {
	body, err := GetEnumReserveInfo("192.168.0.59:5510")
	if err != nil {
		logging.Error("failed to get enum reserve info", slog.Any("error", err))
	}

	var entry Entry
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		logging.Error("failed to unmarshal xml", slog.Any("error", err))
	}

	hasReserve, timeList, err := HasRemainReserve(&entry)
	if hasReserve {
		if err != nil {
			logging.Error("failed to check remaining reserve", slog.Any("error", err))
		}
		logging.InfoAttrs(context.Background(), slog.LevelInfo, "test EDCB API (EnumReserveInfo)", slog.Any("timeList", timeList), slog.Bool("hasReserve", hasReserve))

	}
}
