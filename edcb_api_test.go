package main

import (
	"encoding/xml"
	"log"
	"testing"
)

func TestEDCBAPI(t *testing.T) {
	url := "http://localhost:5510/api/EnumReserveInfo"
	body, err := APIReq2Body(url)
	if err != nil {
		Errorlog(err)
	}

	var entry Entry
	err = xml.Unmarshal(body, &entry)
	if err != nil {
		Errorlog(err)
	}

	hasReserve, timeList, err := HasRemainReserve(&entry)
	if hasReserve {
		if err != nil {
			Errorlog(err)
		}
		log.Println(timeList)
		log.Println(hasReserve)
	}
}
