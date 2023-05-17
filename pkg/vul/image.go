package vul

import (
	"time"
)

type SummaryItem struct {
	Image           string    `json:"image,omitempty"`
	ImageCount      int       `json:"image_count"`
	VersionCount    int       `json:"version_count"`
	SourceCount     int       `json:"source_count"`
	PackageCount    int       `json:"package_count"`
	LastReading     time.Time `json:"last_reading"`
	TotalExposures  int       `json:"total_exposures"`
	UniqueExposures int       `json:"unique_exposures"`
	FixedCount      int       `json:"fixed_count"`
	Exposure        Exposure  `json:"exposure"`

	*Item `json:"-"`
}

type Exposure struct {
	Negligible int `json:"negligible"`
	Low        int `json:"low"`
	Medium     int `json:"medium"`
	High       int `json:"high"`
	Critical   int `json:"critical"`
	Unknown    int `json:"unknown"`

	*Item `json:"-"`
}
