package vul

type ListImageTimelineItem struct {
	Sources map[string]*ListImageSourceTimelineItem `json:"sources"`
}

type ListImageSourceTimelineItem struct {
	Total      int `json:"total"`
	Negligible int `json:"negligible"`
	Low        int `json:"low"`
	Medium     int `json:"medium"`
	High       int `json:"high"`
	Critical   int `json:"critical"`
	Unknown    int `json:"unknown"`
}
