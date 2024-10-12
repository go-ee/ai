package api

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	ContextName      string
	SessionName      string
	PatternName      string
	PatternVariables map[string]string
	Message          string
	Language         string
	Meta             string
}

type ChatOptions struct {
	Model            string
	Temperature      float64
	TopP             float64
	PresencePenalty  float64
	FrequencyPenalty float64
	Raw              bool
	Seed             int
}
