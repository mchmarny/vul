package query

import "time"

type ListImageItem struct {
	Image        string    `json:"image"`
	VersionCount int       `json:"version_count"`
	FirstReading time.Time `json:"first_reading"`
	LastReading  time.Time `json:"last_reading"`
}
