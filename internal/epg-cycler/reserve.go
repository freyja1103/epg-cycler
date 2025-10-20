package epgcycler

import (
	"time"

	edcbapiclient "github.com/freyja1103/epg-cycler/edcb-api-client"
)

type ReserveInfo struct {
	Title     string
	StartDate string
	StartTime string
}

func (e *epgCycler) HasReserve(entry *edcbapiclient.ReserveInfoEntry) (bool, []*ReserveInfo, error) {
	now := time.Now()
	nowDate, err := time.Parse("2006/01/02", now.Format("2006/01/02"))
	if err != nil {
		return false, nil, err
	}
	var timeList []*ReserveInfo

	for _, reserve := range entry.Items.ReserveInfo {
		startDate, err := parseDateBySlash(reserve.StartDate)
		if err != nil {
			return false, nil, err
		}
		startTime, err := parseTimeByColon(reserve.StartTime)
		if err != nil {
			return false, nil, err
		}

		// 同じ日付または次の日の{ReserveCutoffHour}時までの予約を対象とする
		if (0 <= now.Hour() || now.Hour() <= e.Config.ReserveCutoffHour) && (startDate.Day() == now.Day() && startTime.Hour() < e.Config.ReserveCutoffHour) {
			timeList = append(timeList, &ReserveInfo{
				Title:     reserve.Title,
				StartDate: reserve.StartDate,
				StartTime: reserve.StartTime,
			})
		}

		// 次の日の予約
		if nowDate.Day()+1 == startDate.Day() && startTime.Hour() < e.Config.ReserveCutoffHour {
			timeList = append(timeList, &ReserveInfo{
				Title:     reserve.Title,
				StartDate: reserve.StartDate,
				StartTime: reserve.StartTime,
			})
		}
	}
	return len(timeList) >= 1, timeList, nil
}

func parseDateBySlash(d string) (*time.Time, error) {
	stDate, err := time.Parse("2006/01/02", d)
	if err != nil {
		return &time.Time{}, err
	}
	return &stDate, nil
}

func parseTimeByColon(t string) (*time.Time, error) {
	stTime, err := time.Parse("15:04:05", t)
	if err != nil {
		return &time.Time{}, err
	}
	return &stTime, nil
}
