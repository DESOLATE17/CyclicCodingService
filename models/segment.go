package models

import "time"

type Data struct {
	NumOfSegment  int       `json:"numOfSegment"`
	TotalSegments int       `json:"totalSegments"`
	Message       []byte    `json:"message"`
	Date          time.Time `json:"sendDate"`
	Sender        string    `json:"sender"`
	MessageId     int       `json:"messageId"`
}
