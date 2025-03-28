package models

import "time"

// Document 文档模型
type Document struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"size:255;not null"`
	Content    string    `json:"content" gorm:"type:text"`
	Type       string    `json:"type" gorm:"size:50;not null"` // public, private
	CategoryID *uint     `json:"category_id"`
	CreatorID  uint      `json:"creator_id" gorm:"not null"`
	Status     int       `json:"status" gorm:"default:1"` // 1:active, 0:deleted
	Version    int       `json:"version" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DocumentCategory 文档分类
type DocumentCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"` // 父分类ID，支持层级分类
	CreatorID   uint      `json:"creator_id" gorm:"not null"`
	Status      int       `json:"status" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DocumentTag 文档标签
type DocumentTag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null;unique"`
	CreatorID uint      `json:"creator_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DocumentTagRelation 文档-标签关系
type DocumentTagRelation struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DocumentID uint      `json:"document_id" gorm:"not null"`
	TagID      uint      `json:"tag_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

// DocumentShare 文档分享
type DocumentShare struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DocumentID uint      `json:"document_id" gorm:"not null"`
	SharedBy   uint      `json:"shared_by" gorm:"not null"`   // 分享人ID
	SharedWith uint      `json:"shared_with" gorm:"not null"` // 被分享人ID
	Permission string    `json:"permission" gorm:"size:50"`   // read, write, admin
	ExpireAt   time.Time `json:"expire_at"`                   // 分享过期时间
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DocumentVersion 文档版本历史
type DocumentVersion struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DocumentID uint      `json:"document_id" gorm:"not null"`
	Version    int       `json:"version" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text"`
	UpdatedBy  uint      `json:"updated_by" gorm:"not null"`
	Comment    string    `json:"comment"` // 版本说明
	CreatedAt  time.Time `json:"created_at"`
}

// DocumentComment 文档评论
type DocumentComment struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DocumentID uint      `json:"document_id" gorm:"not null"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	ParentID   *uint     `json:"parent_id"`               // 父评论ID，支持评论嵌套
	Status     int       `json:"status" gorm:"default:1"` // 1:active, 0:deleted
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
