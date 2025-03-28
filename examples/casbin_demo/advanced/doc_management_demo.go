package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

// Document 表示一个文档
type Document struct {
	ID      string
	Name    string
	Type    string // "private" 或 "public"
	Content string
	OwnerID string
}

// DocumentManager 文档管理器
type DocumentManager struct {
	enforcer *casbin.Enforcer
	docs     map[string]*Document
}

// NewDocumentManager 创建文档管理器
func NewDocumentManager() (*DocumentManager, error) {
	e, err := casbin.NewEnforcer("doc_management_model.conf", "doc_management_policy.csv")
	if err != nil {
		return nil, fmt.Errorf("创建enforcer失败: %v", err)
	}

	return &DocumentManager{
		enforcer: e,
		docs:     make(map[string]*Document),
	}, nil
}

// AddDocument 添加文档
func (dm *DocumentManager) AddDocument(doc *Document) error {
	dm.docs[doc.ID] = doc
	// 添加文档类型关系
	_, err := dm.enforcer.AddNamedGroupingPolicy("g3", doc.ID, fmt.Sprintf("doc_%s", doc.Type))
	return err
}

// CheckPermission 检查权限
func (dm *DocumentManager) CheckPermission(userID, docID, action string) bool {
	ok, err := dm.enforcer.Enforce(userID, docID, action)
	if err != nil {
		log.Printf("检查权限失败: %v", err)
		return false
	}
	return ok
}

// GetDocument 获取文档（带权限检查）
func (dm *DocumentManager) GetDocument(userID, docID string) (*Document, error) {
	if !dm.CheckPermission(userID, docID, "read") {
		return nil, fmt.Errorf("没有权限读取文档 %s", docID)
	}

	doc, exists := dm.docs[docID]
	if !exists {
		return nil, fmt.Errorf("文档 %s 不存在", docID)
	}
	return doc, nil
}

// UpdateDocument 更新文档（带权限检查）
func (dm *DocumentManager) UpdateDocument(userID string, doc *Document) error {
	if !dm.CheckPermission(userID, doc.ID, "write") {
		return fmt.Errorf("没有权限修改文档 %s", doc.ID)
	}

	dm.docs[doc.ID] = doc
	return nil
}

func ExampleDocumentManagement() {
	dm, err := NewDocumentManager()
	if err != nil {
		log.Fatalf("创建文档管理器失败: %v", err)
	}

	// 创建示例文档
	docs := []*Document{
		{ID: "doc1", Name: "私密报告", Type: "private", Content: "机密内容", OwnerID: "alice"},
		{ID: "doc2", Name: "项目计划", Type: "private", Content: "项目详情", OwnerID: "bob"},
		{ID: "doc3", Name: "公共通知", Type: "public", Content: "公告内容", OwnerID: "alice"},
		{ID: "doc4", Name: "使用手册", Type: "public", Content: "使用说明", OwnerID: "bob"},
	}

	// 添加文档到系统
	for _, doc := range docs {
		if err := dm.AddDocument(doc); err != nil {
			log.Printf("添加文档失败: %v", err)
		}
	}

	// 测试用例
	testCases := []struct {
		userID string
		docID  string
		action string
		desc   string
	}{
		{"alice", "doc1", "read", "管理员读取私有文档"},
		{"alice", "doc1", "write", "管理员修改私有文档"},
		{"bob", "doc1", "read", "组长读取私有文档"},
		{"bob", "doc2", "write", "组长修改私有文档"},
		{"charles", "doc3", "read", "普通用户读取公共文档"},
		{"charles", "doc1", "read", "普通用户尝试读取私有文档"},
		{"david", "doc4", "read", "普通用户读取公共文档"},
		{"eve", "doc2", "write", "普通用户尝试修改私有文档"},
	}

	fmt.Println("\n=== 文档管理系统权限测试 ===")
	for _, tc := range testCases {
		ok := dm.CheckPermission(tc.userID, tc.docID, tc.action)
		fmt.Printf("用户: %-8s | 操作: %-6s | 文档: %-6s | %-30s | 允许访问: %v\n",
			tc.userID,
			tc.action,
			tc.docID,
			tc.desc,
			ok)
	}

	// 演示文档访问
	fmt.Println("\n=== 文档访问测试 ===")
	testAccesses := []struct {
		userID string
		docID  string
		desc   string
	}{
		{"alice", "doc1", "管理员访问私有文档"},
		{"bob", "doc2", "组长访问私有文档"},
		{"charles", "doc3", "普通用户访问公共文档"},
		{"charles", "doc1", "普通用户尝试访问私有文档"},
	}

	for _, test := range testAccesses {
		doc, err := dm.GetDocument(test.userID, test.docID)
		if err != nil {
			fmt.Printf("%s: %v\n", test.desc, err)
			continue
		}
		fmt.Printf("%s: 成功访问文档 '%s'\n", test.desc, doc.Name)
	}
}
