package api

import (
	"chatgpt-ai/api/chat"

	"github.com/gin-gonic/gin"
)

type ApiInerface interface {
	RegisterRouter(engine *gin.Engine)
}

var ApiList []ApiInerface

func init() {
	ApiList = append(ApiList, &chat.ChatApiRouter{})
}

func RegisterRouter(engine *gin.Engine) {
	for _, api := range ApiList {
		api.RegisterRouter(engine)
	}
}
