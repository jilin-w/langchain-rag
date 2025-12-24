package chat_service

import (
	"chatgpt-ai/model/chat_model"
	"chatgpt-ai/service/ollama_service"

	"github.com/gin-gonic/gin"
)

var DefaultChatService = &chatService{}

type chatService struct {
}

func (*chatService) Chat(ctx *gin.Context) {
	var req chat_model.ChatReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}
	resp, err := ollama_service.DefaultOllamaService.Chat(req.Prompt, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, chat_model.ChatResp{Response: resp})
}

func (*chatService) ChatStream(ctx *gin.Context) {
	var req chat_model.ChatReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}
	ollama_service.DefaultOllamaService.ChatStream(req.Prompt, nil, ctx)
}
