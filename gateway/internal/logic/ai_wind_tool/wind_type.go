package ai_wind_tool

type StreamEvent struct {
	Type       string            `json:"type"`
	Content    string            `json:"content,omitempty"`
	Answer     string            `json:"answer,omitempty"`
	TraceID    string            `json:"traceId,omitempty"`
	Citations  []StreamCitation  `json:"citations,omitempty"`
	References []StreamReference `json:"references,omitempty"`
}

type StreamReference struct {
	RefIndex int            `json:"refIndex"`
	RefText  string         `json:"refText"`
	Start    int            `json:"start"`
	End      int            `json:"end"`
	Citation StreamCitation `json:"citation"`
}

type StreamCitation struct {
	RefIndex   int     `json:"refIndex"`
	RefText    string  `json:"refText"`
	DocumentId int64   `json:"documentId"`
	ChunkId    int64   `json:"chunkId"`
	Title      string  `json:"title"`
	Snippet    string  `json:"snippet"`
	Score      float64 `json:"score"`
}
