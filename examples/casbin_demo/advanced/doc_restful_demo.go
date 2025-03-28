package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Document 表示一个文档
type RESTDocument struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"` // "private" 或 "public"
	Content string `json:"content"`
	OwnerID string `json:"owner_id"`
}

// Comment 表示文档评论
type Comment struct {
	ID        string `json:"id"`
	DocID     string `json:"doc_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// DocumentAPI RESTful API 处理器
type DocumentAPI struct {
	enforcer *casbin.Enforcer
	docs     map[string]*RESTDocument
	comments map[string][]*Comment
}

// NewDocumentAPI 创建文档 API 处理器
func NewDocumentAPI() (*DocumentAPI, error) {
	e, err := casbin.NewEnforcer("doc_restful_model.conf", "doc_restful_policy.csv")
	if err != nil {
		return nil, fmt.Errorf("创建enforcer失败: %v", err)
	}

	return &DocumentAPI{
		enforcer: e,
		docs:     make(map[string]*RESTDocument),
		comments: make(map[string][]*Comment),
	}, nil
}

// AuthMiddleware 权限检查中间件
func (api *DocumentAPI) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader("X-User-ID") // 从请求头获取用户ID
		if user == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 构建权限检查请求
		path := c.Request.URL.Path
		method := c.Request.Method

		// 检查权限
		ok, err := api.enforcer.Enforce(user, path, method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// setupRouter 设置路由
func (api *DocumentAPI) setupRouter() *gin.Engine {
	r := gin.Default()

	// 使用认证中间件
	r.Use(api.AuthMiddleware())

	v1 := r.Group("/api/v1")
	{
		// 文档相关接口
		docs := v1.Group("/documents")
		{
			// 公共文档接口
			public := docs.Group("/public")
			{
				public.GET("", api.ListPublicDocuments)
				public.GET("/:id", api.GetDocument)
				public.POST("", api.CreateDocument)
				public.PUT("/:id", api.UpdateDocument)
				public.DELETE("/:id", api.DeleteDocument)

				// 评论相关接口
				public.GET("/:id/comments", api.ListComments)
				public.POST("/:id/comments", api.CreateComment)
			}

			// 私有文档接口
			private := docs.Group("/private")
			{
				private.GET("", api.ListPrivateDocuments)
				private.GET("/:id", api.GetDocument)
				private.POST("", api.CreateDocument)
				private.PUT("/:id", api.UpdateDocument)
				private.DELETE("/:id", api.DeleteDocument)
			}
		}
	}

	return r
}

// API 处理函数
func (api *DocumentAPI) ListPublicDocuments(c *gin.Context) {
	var docs []*RESTDocument
	for _, doc := range api.docs {
		if doc.Type == "public" {
			docs = append(docs, doc)
		}
	}
	c.JSON(http.StatusOK, docs)
}

func (api *DocumentAPI) ListPrivateDocuments(c *gin.Context) {
	var docs []*RESTDocument
	for _, doc := range api.docs {
		if doc.Type == "private" {
			docs = append(docs, doc)
		}
	}
	c.JSON(http.StatusOK, docs)
}

func (api *DocumentAPI) GetDocument(c *gin.Context) {
	id := c.Param("id")
	doc, exists := api.docs[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在"})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (api *DocumentAPI) CreateDocument(c *gin.Context) {
	var doc RESTDocument
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置文档类型路径
	docType := "public"
	if strings.Contains(c.Request.URL.Path, "/private/") {
		docType = "private"
	}
	doc.Type = docType

	api.docs[doc.ID] = &doc
	c.JSON(http.StatusCreated, doc)
}

func (api *DocumentAPI) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	doc, exists := api.docs[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在"})
		return
	}

	var updatedDoc RESTDocument
	if err := c.ShouldBindJSON(&updatedDoc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedDoc.ID = id
	updatedDoc.Type = doc.Type // 保持原有类型
	api.docs[id] = &updatedDoc

	c.JSON(http.StatusOK, updatedDoc)
}

func (api *DocumentAPI) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if _, exists := api.docs[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在"})
		return
	}

	delete(api.docs, id)
	c.Status(http.StatusNoContent)
}

func (api *DocumentAPI) ListComments(c *gin.Context) {
	docID := c.Param("id")
	comments := api.comments[docID]
	c.JSON(http.StatusOK, comments)
}

func (api *DocumentAPI) CreateComment(c *gin.Context) {
	docID := c.Param("id")
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.DocID = docID
	comment.UserID = c.GetHeader("X-User-ID")

	if api.comments[docID] == nil {
		api.comments[docID] = make([]*Comment, 0)
	}
	api.comments[docID] = append(api.comments[docID], &comment)

	c.JSON(http.StatusCreated, comment)
}

func ExampleRESTfulAPI() {
	api, err := NewDocumentAPI()
	if err != nil {
		log.Fatalf("创建API失败: %v", err)
	}

	r := api.setupRouter()
	fmt.Println("启动服务器在 :8080...")
	r.Run(":8080")
}
