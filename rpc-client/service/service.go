package service

import (
	"errors"
	"log"
)

// 用户服务接口
type UserService interface {
	GetUser(id int, user *User) error
	CreateUser(user *User, response *bool) error
}

// 用户结构体
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 用户服务实现
type UserServiceImpl struct {
	Users map[int]*User
}

var ErrUserNotFound = errors.New("user not found")

// 获取用户
func (s *UserServiceImpl) GetUser(id int, user *User) error {
	if u, ok := s.Users[id]; ok {
		*user = *u
		log.Printf("RPC client Get user: %v", user)
		return nil
	}
	return ErrUserNotFound
}

// 创建用户
func (s *UserServiceImpl) CreateUser(user *User, response *bool) error {
	s.Users[user.Id] = user
	log.Printf("create user: %v", user)
	*response = true
	return nil
}
