package vul

type ImageTimeline struct {
	Date  string `json:"date"`
	Grype int    `json:"grype"`
	Trivy int    `json:"trivy"`
	Snyk  int    `json:"snyk"`
}
