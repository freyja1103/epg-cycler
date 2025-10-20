package edcbapiclient

import "encoding/xml"

type ReserveInfoEntry struct {
	XMLName xml.Name `xml:"entry"`
	Total   int      `xml:"total"`
	Index   int      `xml:"index"`
	Count   int      `xml:"count"`
	Items   Items    `xml:"items"`
}
