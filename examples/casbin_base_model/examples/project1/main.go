package main

import (
	"fmt"
	"log"

	"casbin_base_model/enforcer"
)

func main() {
	// 创建配置
	config := &enforcer.Config{
		DBType:       "mysql",
		DBConnection: "root:password@tcp(127.0.0.1:3306)/casbin_demo?charset=utf8mb4&parseTime=True&loc=Local",
		ModelPath:    "../../models/rbac_with_domains.conf",
		AutoLoad:     true,
	}

	// 创建 enforcer
	e, err := enforcer.NewEnforcer(config)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}

	// 添加策略
	// 在 project1 域中定义权限
	policies := [][]interface{}{
		// 管理员权限
		{"admin", "project1", "/api/*", "*"},
		// 开发者权限
		{"developer", "project1", "/api/v1/products/*", "GET"},
		{"developer", "project1", "/api/v1/products/*", "POST"},
		// 用户权限
		{"user", "project1", "/api/v1/products/*", "GET"},
	}

	// 添加角色分配
	roles := [][]interface{}{
		{"alice", "admin", "project1"},
		{"bob", "developer", "project1"},
		{"charles", "user", "project1"},
	}

	// 批量添加策略
	for _, policy := range policies {
		if _, err := e.AddPolicy(policy...); err != nil {
			log.Printf("Failed to add policy %v: %v", policy, err)
		}
	}

	// 批量添加角色
	for _, role := range roles {
		if _, err := e.AddGroupingPolicy(role...); err != nil {
			log.Printf("Failed to add role %v: %v", role, err)
		}
	}

	// 测试权限
	testCases := []struct {
		user   string
		domain string
		path   string
		method string
	}{
		{"alice", "project1", "/api/v1/products", "GET"},
		{"alice", "project1", "/api/v2/settings", "POST"},
		{"bob", "project1", "/api/v1/products", "GET"},
		{"bob", "project1", "/api/v2/settings", "GET"},
		{"charles", "project1", "/api/v1/products", "GET"},
		{"charles", "project1", "/api/v1/products", "POST"},
	}

	fmt.Println("\n权限测试结果:")
	for _, tc := range testCases {
		ok, err := e.Enforce(tc.user, tc.domain, tc.path, tc.method)
		if err != nil {
			log.Printf("Enforce failed: %v", err)
			continue
		}
		fmt.Printf("用户: %-8s | 域: %-8s | 路径: %-20s | 方法: %-6s | 允许访问: %v\n",
			tc.user, tc.domain, tc.path, tc.method, ok)
	}
}
