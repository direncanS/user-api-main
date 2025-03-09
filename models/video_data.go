package models

import "time"

type VideoData struct {
	Id        string    `json:"id"`
	Folder    string    `json:"folder"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
