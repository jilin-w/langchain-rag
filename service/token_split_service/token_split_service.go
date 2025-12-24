package token_split_service

import (
	"context"
	"io"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

var DeFaultTokenSplit = &TokenSplitService{}

type TokenSplitService struct {
}

// 使用LangChainTextSplitter进行分词
func (*TokenSplitService) Split(reader io.Reader) ([]schema.Document, error) {
	loader := documentloaders.NewText(reader)
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 128
	split.ChunkOverlap = 30
	docs, err := loader.LoadAndSplit(context.Background(), split)
	return docs, err
}
