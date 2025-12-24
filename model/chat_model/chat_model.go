package chat_model

type ChatReq struct {
	Prompt string `json:"prompt" form:"prompt"`
}
type ChatResp struct {
	Response string `json:"response"`
}

type ChatRagReq struct {
	Prompt  string `json:"prompt" form:"prompt"`
	RagName string `json:"rag_name" form:"rag_name"`
}
