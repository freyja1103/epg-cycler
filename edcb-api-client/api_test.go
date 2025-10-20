package edcbapiclient

import (
	"context"
	"log/slog"
	"testing"

	"github.com/freyja1103/epg-cycler/logging"
)

func TestEDCBAPI(t *testing.T) {

	edcbapi := NewEDCBAPIClient("192.168.0.59:5510", nil)
	_, err := edcbapi.GetEnumReserveInfo(context.Background())
	if err != nil {
		logging.Error("failed to get enum reserve info", slog.Any("error", err))
	}

	// hasReserve, timeList, err := HasRemainReserve(&entry)
	// if hasReserve {
	// 	if err != nil {
	// 		logging.Error("failed to check remaining reserve", slog.Any("error", err))
	// 	}
	// 	logging.InfoAttrs(context.Background(), slog.LevelInfo, "test EDCB API (EnumReserveInfo)", slog.Any("timeList", timeList), slog.Bool("hasReserve", hasReserve))

	// }
}
