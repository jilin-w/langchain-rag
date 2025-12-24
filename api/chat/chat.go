package chat

import (
	"chatgpt-ai/service/chat_service"
	"chatgpt-ai/service/ollama_service"

	"github.com/gin-gonic/gin"
)

type ChatApiRouter struct {
}

func (*ChatApiRouter) RegisterRouter(engine *gin.Engine) {
	group := engine.Group("/api/gpt")
	group.POST("/chat", ChatApi)
	group.GET("/chat/stream", ChatStreamApi)
	group.POST("/rag/upload", UploadRagFile)
	group.GET("/rag/chat", ChatWithRag)
}
func ChatApi(ctx *gin.Context) {
	chat_service.DefaultChatService.Chat(ctx)
}

func ChatStreamApi(ctx *gin.Context) {
	chat_service.DefaultChatService.ChatStream(ctx)
}

func UploadRagFile(ctx *gin.Context) {
	ollama_service.DefaultOllamaService.UploadRagFile(ctx)
}

func ChatWithRag(ctx *gin.Context) {
	ollama_service.DefaultOllamaService.ChatWithRag(ctx)
}
