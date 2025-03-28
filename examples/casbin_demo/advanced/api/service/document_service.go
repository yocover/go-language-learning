package service

import (
	"fmt"
	"time"

	"github.com/go-language-learning/examples/casbin_demo/advanced/api/models"
	"gorm.io/gorm"
)

// DocumentService 文档服务
type DocumentService struct {
	db          *gorm.DB
	authService *AuthService
}

// NewDocumentService 创建文档服务
func NewDocumentService(db *gorm.DB, authService *AuthService) *DocumentService {
	return &DocumentService{
		db:          db,
		authService: authService,
	}
}

// CreateDocument 创建文档
func (s *DocumentService) CreateDocument(userID uint, doc *models.Document) error {
	// 检查用户是否有创建文档的权限
	allowed, err := s.authService.CheckPermission(userID, "platform", "/api/documents", "POST")
	if err != nil {
		return fmt.Errorf("检查权限失败: %w", err)
	}
	if !allowed {
		return fmt.Errorf("没有创建文档的权限")
	}

	// 设置创建者ID
	doc.CreatorID = userID
	doc.Version = 1

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建文档
		if err := tx.Create(doc).Error; err != nil {
			return fmt.Errorf("创建文档失败: %w", err)
		}

		// 创建首个版本记录
		version := &models.DocumentVersion{
			DocumentID: doc.ID,
			Version:    1,
			Content:    doc.Content,
			UpdatedBy:  userID,
			Comment:    "初始版本",
		}
		if err := tx.Create(version).Error; err != nil {
			return fmt.Errorf("创建版本记录失败: %w", err)
		}

		return nil
	})
}

// UpdateDocument 更新文档
func (s *DocumentService) UpdateDocument(userID uint, docID uint, updates map[string]interface{}) error {
	// 检查文档是否存在
	var doc models.Document
	if err := s.db.First(&doc, docID).Error; err != nil {
		return fmt.Errorf("文档不存在: %w", err)
	}

	// 检查权限
	allowed, err := s.checkDocumentAccess(userID, docID, "PUT")
	if err != nil {
		return err
	}
	if !allowed {
		return fmt.Errorf("没有更新文档的权限")
	}

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新版本号
		updates["version"] = doc.Version + 1

		// 更新文档
		if err := tx.Model(&doc).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新文档失败: %w", err)
		}

		// 创建新版本记录
		version := &models.DocumentVersion{
			DocumentID: docID,
			Version:    doc.Version + 1,
			Content:    updates["content"].(string),
			UpdatedBy:  userID,
			Comment:    updates["comment"].(string),
		}
		if err := tx.Create(version).Error; err != nil {
			return fmt.Errorf("创建版本记录失败: %w", err)
		}

		return nil
	})
}

// ShareDocument 分享文档
func (s *DocumentService) ShareDocument(userID uint, share *models.DocumentShare) error {
	// 检查文档是否存在
	var doc models.Document
	if err := s.db.First(&doc, share.DocumentID).Error; err != nil {
		return fmt.Errorf("文档不存在: %w", err)
	}

	// 检查分享权限
	allowed, err := s.checkDocumentAccess(userID, share.DocumentID, "POST")
	if err != nil {
		return err
	}
	if !allowed {
		return fmt.Errorf("没有分享文档的权限")
	}

	// 设置分享人ID
	share.SharedBy = userID

	// 创建分享记录
	if err := s.db.Create(share).Error; err != nil {
		return fmt.Errorf("创建分享记录失败: %w", err)
	}

	return nil
}

// AddComment 添加评论
func (s *DocumentService) AddComment(comment *models.DocumentComment) error {
	// 检查文档是否存在
	var doc models.Document
	if err := s.db.First(&doc, comment.DocumentID).Error; err != nil {
		return fmt.Errorf("文档不存在: %w", err)
	}

	// 检查评论权限
	allowed, err := s.authService.CheckPermission(comment.UserID, "platform",
		fmt.Sprintf("/api/documents/%d/comments", comment.DocumentID), "POST")
	if err != nil {
		return fmt.Errorf("检查权限失败: %w", err)
	}
	if !allowed {
		return fmt.Errorf("没有评论权限")
	}

	// 创建评论
	if err := s.db.Create(comment).Error; err != nil {
		return fmt.Errorf("创建评论失败: %w", err)
	}

	return nil
}

// GetDocument 获取文档
func (s *DocumentService) GetDocument(userID uint, docID uint) (*models.Document, error) {
	var doc models.Document
	if err := s.db.First(&doc, docID).Error; err != nil {
		return nil, fmt.Errorf("文档不存在: %w", err)
	}

	// 检查访问权限
	allowed, err := s.checkDocumentAccess(userID, docID, "GET")
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("没有访问文档的权限")
	}

	return &doc, nil
}

// GetDocumentVersions 获取文档版本历史
func (s *DocumentService) GetDocumentVersions(userID uint, docID uint) ([]models.DocumentVersion, error) {
	// 检查访问权限
	allowed, err := s.checkDocumentAccess(userID, docID, "GET")
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("没有访问文档的权限")
	}

	var versions []models.DocumentVersion
	if err := s.db.Where("document_id = ?", docID).Order("version desc").Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("获取版本历史失败: %w", err)
	}

	return versions, nil
}

// GetDocumentComments 获取文档评论
func (s *DocumentService) GetDocumentComments(userID uint, docID uint) ([]models.DocumentComment, error) {
	// 检查访问权限
	allowed, err := s.checkDocumentAccess(userID, docID, "GET")
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("没有访问文档的权限")
	}

	var comments []models.DocumentComment
	if err := s.db.Where("document_id = ? AND status = 1", docID).
		Order("created_at desc").Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("获取评论失败: %w", err)
	}

	return comments, nil
}

// checkDocumentAccess 检查文档访问权限
func (s *DocumentService) checkDocumentAccess(userID uint, docID uint, action string) (bool, error) {
	var doc models.Document
	if err := s.db.First(&doc, docID).Error; err != nil {
		return false, fmt.Errorf("文档不存在: %w", err)
	}

	// 检查是否是文档创建者
	if doc.CreatorID == userID {
		return true, nil
	}

	// 检查是否有分享记录
	var share models.DocumentShare
	err := s.db.Where("document_id = ? AND shared_with = ? AND expire_at > ?",
		docID, userID, time.Now()).First(&share).Error
	if err == nil {
		// 检查分享权限
		switch action {
		case "GET":
			return true, nil
		case "PUT":
			return share.Permission == "write" || share.Permission == "admin", nil
		case "DELETE":
			return share.Permission == "admin", nil
		}
	}

	// 检查通用权限
	path := fmt.Sprintf("/api/documents/%d", docID)
	return s.authService.CheckPermission(userID, "platform", path, action)
}
