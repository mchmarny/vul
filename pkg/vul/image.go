package vul

import "time"

type SummaryItem struct {
	Image        string    `json:"image"`
	ImageCount   int       `json:"image_count"`
	VersionCount int       `json:"version_count"`
	SourceCount  int       `json:"source_count"`
	PackageCount int       `json:"package_count"`
	Exposure     Exposure  `json:"exposure"`
	FirstReading time.Time `json:"first_reading"`
	LastReading  time.Time `json:"last_reading"`
}

type Exposure struct {
	Total      int `json:"total"`
	Negligible int `json:"negligible"`
	Low        int `json:"low"`
	Medium     int `json:"medium"`
	High       int `json:"high"`
	Critical   int `json:"critical"`
	Unknown    int `json:"unknown"`
}
