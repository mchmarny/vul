package vul

import "time"

type ListImageVersionRequest struct {
	Image string `json:"image"`
}

type ListImageSourceItem struct {
	Source       string    `json:"source"`
	PackageCount int       `json:"package_count"`
	FirstReading time.Time `json:"first_reading"`
	LastReading  time.Time `json:"last_reading"`
}
