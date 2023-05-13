package vul

type ImageDigestExposures struct {
	Image    string                       `json:"image"`
	Digest   string                       `json:"digest"`
	Packages map[string]*PackageExposures `json:"packages"`
}

type PackageExposures struct {
	Versions map[string]*PackageVersionExposures `json:"versions"`
}

type PackageVersionExposures struct {
	Sources map[string]*SourceExposures `json:"sources"`
}

type SourceExposures struct {
	Exposures map[string]*Exposures `json:"exposures"`
}

type Exposures struct {
	Severity string  `json:"severity"`
	Score    float64 `json:"score"`
	Fixed    bool    `json:"fixed"`
}
