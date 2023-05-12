package vul

type ImageTimeline struct {
	Date  string `json:"date"`
	Name  string `json:"source"`
	Value int    `json:"total"`
}
