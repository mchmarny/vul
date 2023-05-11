package query

type ListImageTimelineRequest struct {
	Image   string `json:"image"`
	FromDay string `json:"from_day,omitempty"`
	ToDay   string `json:"to_day,omitempty"`
}

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
