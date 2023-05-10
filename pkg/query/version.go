package query

import "time"

type ListImageVersionRequest struct {
	Image string `json:"image"`
}

type ListImageVersionItem struct {
	Digest       string    `json:"digest"`
	SourceCount  int       `json:"source_count"`
	FirstReading time.Time `json:"first_reading"`
	LastReading  time.Time `json:"last_reading"`
	PackageCount int       `json:"package_count"`
}
