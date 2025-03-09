package models

import "time"

type WebcamData struct {
	Id       string            `json:"id"`
	Topic    string            `json:"topic"`
	CreateAt time.Time         `json:"created_at"`
	Metadata map[string]string `json:"metadata"`
}
