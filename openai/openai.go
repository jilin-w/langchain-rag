package openai

import (
	"github.com/tmc/langchaingo/llms/openai"
)

var lLMService *openai.LLM

func init() {
	llm, err := openai.New(openai.WithBaseURL("http://127.0.0.1:11434/v1"), openai.WithModel("deepseek-r1:8b"), openai.WithEmbeddingModel("nomic-embed-text:v1.5"))
	if err != nil {
		panic(err)
	}
	lLMService = llm
}

func GetLLMInstance() *openai.LLM {
	return lLMService
}
