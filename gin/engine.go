package gin

import (
	"chatgpt-ai/api"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	//实现自动注册路由
	api.RegisterRouter(r)
	r.Run(":3001") // listen and serve on 0.0.0.0:8080
}
