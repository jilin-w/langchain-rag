package token_persistence_service

import (
	"chatgpt-ai/openai"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

var defaultPersistenceService PersistenceService

func init() {
	defaultPersistenceService = &PostgresPersistenceService{}
}
func Persistence(docs []schema.Document) (err error) {
	return defaultPersistenceService.persistence(docs)
}
func Search(name, query string) (docs []schema.Document, err error) {
	return defaultPersistenceService.search(name, query)
}

// 持久化接口
type PersistenceService interface {
	persistence(docs []schema.Document) (err error)
	search(name, query string) (docs []schema.Document, err error)
}

type PostgresPersistenceService struct {
}

func (*PostgresPersistenceService) persistence(docs []schema.Document) (err error) {
	llm := openai.GetLLMInstance()
	embedder, err := embeddings.NewEmbedder(llm)
	store, err := pgvector.New(context.Background(),
		pgvector.WithConnectionURL("postgresql://postgresql:123456@127.0.0.1:5432/postgres?sslmode=disable"),
		pgvector.WithEmbedder(embedder),
	)
	if err != nil {
		return
	}
	_, err = store.AddDocuments(context.Background(), docs, vectorstores.WithEmbedder(embedder))
	if err != nil {
		fmt.Println("Error adding documents:", err)
		return
	}
	return
}

func (*PostgresPersistenceService) search(name, query string) (docs []schema.Document, err error) {
	llm := openai.GetLLMInstance()
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return
	}
	store, err := pgvector.New(context.Background(),
		pgvector.WithConnectionURL("postgresql://postgresql:123456@127.0.0.1:5432/postgres?sslmode=disable"),
		pgvector.WithEmbedder(embedder),
	)
	if err != nil {
		return
	}
	docs, err = store.SimilaritySearch(context.Background(), query, 50,
		vectorstores.WithEmbedder(embedder))
	return
}
