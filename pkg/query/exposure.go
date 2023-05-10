package query

type ListImageVersionExposureRequest struct {
	Image  string `json:"image"`
	Digest string `json:"digest"`
}

type ListDigestExposureItem struct {
	Source   string  `json:"source"`
	Package  string  `json:"package"`
	Version  string  `json:"version"`
	Severity string  `json:"severity"`
	Score    float64 `json:"score"`
	Fixed    bool    `json:"fixed"`
}
