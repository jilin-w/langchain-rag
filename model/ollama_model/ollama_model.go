package ollama_model

type OllamaChatReq struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options"`
}

type OllamaChatResp struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}
