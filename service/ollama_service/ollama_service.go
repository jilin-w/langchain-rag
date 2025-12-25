package ollama_service

import (
	"chatgpt-ai/model/chat_model"
	"chatgpt-ai/openai"
	"chatgpt-ai/service/git_service"
	"chatgpt-ai/service/tika_service"
	"chatgpt-ai/service/token_persistence_service"
	"chatgpt-ai/service/token_split_service"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

var DefaultOllamaService ollamaService

type ollamaService struct {
}

func (s *ollamaService) Chat(prompt string, options map[string]interface{}) (response string, err error) {
	client := openai.GetLLMInstance()
	chatCompletion, err := llms.GenerateFromSinglePrompt(context.Background(), client, prompt, llms.WithModel("deepseek-r1:8b"))
	if err != nil {
		return
	}
	return chatCompletion, nil
}

func (s *ollamaService) ChatStream(prompt string, options map[string]interface{}, ctx *gin.Context) {

	client := openai.GetLLMInstance()
	//回调
	var streamCallback = func(c context.Context, chunk []byte) (err error) {
		ctx.SSEvent("message", map[string]interface{}{
			"content": string(chunk),
		})
		return nil
	}
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	ctx.Writer.Header().Set("Content-Type", "text/plain")
	llms.GenerateFromSinglePrompt(context.Background(), client, prompt,
		llms.WithStreamingFunc(streamCallback),
		llms.WithModel("deepseek-r1:8b"),
	)
	return
}

func (s *ollamaService) UploadRagFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	fileName := file.Filename
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "file is required",
		})
		return
	}
	fileBuf, err := file.Open()
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "file parse faild " + err.Error(),
		})
		return
	}
	reader, err := tika_service.ReadFile(fileBuf)
	defer fileBuf.Close()
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "file parse faild " + err.Error(),
		})
		return
	}
	//分割
	docs, err := token_split_service.DeFaultTokenSplit.Split(reader)
	defer reader.Close()
	for _, doc := range docs {
		doc.Metadata["knowledge"] = fileName
	}
	//存储
	token_persistence_service.Persistence(fileName, docs)
	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    "",
	})
}

const SystemPrompt = `根据上下文，回答以下问题:{{.question}}`

func (s *ollamaService) ChatWithRag(ctx *gin.Context) {
	var req chat_model.ChatRagReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}
	//先查询rag
	docs, err := token_persistence_service.Search(req.RagName, req.Prompt)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}
	llm := openai.GetLLMInstance()
	//系统提示词格式化
	history := memory.NewChatMessageHistory()
	for _, v := range docs {
		history.AddAIMessage(ctx, v.PageContent)
	}
	// 使用历史记录创建一个新的对话缓冲区
	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))
	var streamCallback = func(c context.Context, chunk []byte) (err error) {
		ctx.SSEvent("message", map[string]interface{}{
			"content": string(chunk),
		})
		return nil
	}
	chain := chains.NewConversation(llm, conversation)
	res, err := chains.Run(context.Background(), chain, req.Prompt, chains.WithStreamingFunc(streamCallback))
	if err != nil {
		fmt.Printf("error: %v  res %v \n", err, res)
	}
	return
}

// 分析git代码
func (s *ollamaService) AnalyzeGitCode(ctx *gin.Context) {
	repoUrl := ctx.Query("repo_url")
	filePath, err := git_service.CloneRepo(repoUrl)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//先判断文件夹
	dir, _ := os.Stat(filePath)
	if !dir.IsDir() {
		ctx.JSON(400, gin.H{"error": "不是文件夹"})
		return
	}
	entrys, err := os.ReadDir(filePath)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = readFile(strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1], entrys, filePath, &map[string]any{"url": repoUrl, "type": "git_code"})
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{})
}

// 异步写入
func analyzeSigleCodeFile(name string, fileName string, file io.ReadCloser, docs []schema.Document, meta *map[string]any) {
	fmt.Printf("start analyze code file %s \n", fileName)
	reader, err := tika_service.ReadFile(file)
	defer file.Close()
	if err != nil {
		fmt.Printf("tika file error %v \n", err)
		return
	}
	//拆分
	docs, err = token_split_service.DeFaultTokenSplit.Split(reader)
	defer reader.Close()
	if err != nil {
		fmt.Printf("split file error %v \n", err)
		return
	}
	//存储
	for index, doc := range docs {
		doc.Metadata = *meta
		docs[index] = doc
	}
	//存储
	token_persistence_service.Persistence(name, docs)
	fmt.Printf("end analyze code file %s \n", fileName)
}
func readFile(name string, entrys []os.DirEntry, filePath string, meta *map[string]any) (err error) {
	for _, v := range entrys {
		if v.IsDir() {
			dirName := v.Name()
			if strings.HasPrefix(dirName, ".git") || strings.HasPrefix(dirName, ".idea") {
				continue
			}
			dirs, err := os.ReadDir(filePath + "/" + dirName)
			if err != nil {
				return err
			}
			err = readFile(name, dirs, filePath+"/"+v.Name(), meta)
			if err != nil {
				return err
			}
		} else {
			//file的关闭交由使用方关闭
			file, err := os.Open(filePath + "/" + v.Name())
			if err != nil {
				return err
			}
			//tika解析获取
			reader, err := tika_service.ReadFile(file)
			//大模型写入知识库
			go analyzeSigleCodeFile(name, file.Name(), reader, nil, meta)
		}
	}
	return nil
}
