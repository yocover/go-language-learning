package main

// 简化的RAG 系统示例

type Document struct {
	Content string
	Source  string
}

// 检索系统接口
type Retriever interface {
	Search(query string) []Document
}

// 生成系统接口
type Generator interface {
	Generate(context []Document, query string) string
}

// 简单的RAG 系统
type RAGSystem struct {
	retriever Retriever
	generator Generator
}

func (r *RAGSystem) Answer(query string) string {

	// 检索相关文档
	docs := r.retriever.Search(query)

	// 将检索结果和问题一起发送给生成器
	return r.generator.Generate((docs), query)
}
