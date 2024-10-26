package repository

import (
	"encoding/xml"
	"log"
	"time"

	edcbapi "github.com/freyja1103/epg-cycler/infrastructure/edcb-api-client"
	"github.com/freyja1103/epg-cycler/model"
)

type IReservedProgramRepository struct {
	rProgram model.ReservedProgramRepository
}

func NewReservedProgramRepository(rProgram model.ReservedProgramRepository) *IReservedProgramRepository {
	return &IReservedProgramRepository{rProgram: rProgram}
}

func (rp *IReservedProgramRepository) Get(url string) *model.Entry {
	return parseFromEdcb(url)
}

func (rp *IReservedProgramRepository) HasRemainReserve(rsrv model.ReserveInfo) bool {
	return hasReserve(rsrv)
}

func parseFromEdcb(url string) *model.Entry {
	var entry model.Entry
	res, err := edcbapi.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	err = xml.Unmarshal(res, &entry)
	if err != nil {
		log.Fatalln(err)
	}

	return &entry
}

func hasReserve(reserveItem model.ReserveInfo) bool {
	var startDate, startTime time.Time
	var timeList []model.TimeInfo

	now := time.Now()
	date := now.Format("2006/01/02")

	if 0 <= now.Hour() || now.Hour() <= 4 {
		if startDate.Day() == now.Day() && startTime.Hour() < 4 {
			timeList = append(timeList, model.TimeInfo{
				StartDate: reserveItem.StartDate,
				StartTime: reserveItem.StartTime,
			})
		}
	}
	if now.Hour() > 4 && date == reserveItem.StartDate {
		timeList = append(timeList, model.TimeInfo{
			StartDate: reserveItem.StartDate,
			StartTime: reserveItem.StartTime,
		})
	}

	return len(timeList) > 1
}
