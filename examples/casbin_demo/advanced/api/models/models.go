package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey"`
	Username  string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;unique"`
	Status    int       `gorm:"default:1"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// UserGroup 用户组模型
type UserGroup struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"size:255;not null;unique"`
	Description string    `gorm:"size:1000"`
	ParentID    *uint     `gorm:"default:null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// UserGroupMember 用户组成员关系
type UserGroupMember struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	GroupID   uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

// Role 角色模型
type Role struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"size:255;not null;unique"`
	Description string    `gorm:"size:1000"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// UserRole 用户角色关系
type UserRole struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	RoleID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
