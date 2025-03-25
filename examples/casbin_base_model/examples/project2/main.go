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

	// 在 project2 域中定义权限
	// 这个项目使用不同的权限结构
	policies := [][]interface{}{
		// 超级管理员权限
		{"super_admin", "project2", "/*", "*"},
		// 部门管理员权限
		{"dept_admin", "project2", "/department/:id/*", "*"},
		// 员工权限
		{"employee", "project2", "/department/:id/tasks/*", "GET"},
		{"employee", "project2", "/department/:id/tasks/:task_id", "PUT"},
		// 访客权限
		{"guest", "project2", "/public/*", "GET"},
	}

	// 添加角色分配
	roles := [][]interface{}{
		{"david", "super_admin", "project2"},
		{"emma", "dept_admin", "project2"},
		{"frank", "employee", "project2"},
		{"guest1", "guest", "project2"},
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
		desc   string
	}{
		{"david", "project2", "/system/settings", "POST", "超级管理员修改系统设置"},
		{"emma", "project2", "/department/123/settings", "PUT", "部门管理员修改部门设置"},
		{"emma", "project2", "/department/456/settings", "PUT", "部门管理员修改其他部门设置"},
		{"frank", "project2", "/department/123/tasks/list", "GET", "员工查看任务列表"},
		{"frank", "project2", "/department/123/tasks/456", "PUT", "员工更新任务状态"},
		{"frank", "project2", "/department/123/settings", "POST", "员工尝试修改部门设置"},
		{"guest1", "project2", "/public/announcements", "GET", "访客查看公告"},
		{"guest1", "project2", "/department/123/tasks/list", "GET", "访客尝试查看任务列表"},
	}

	fmt.Println("\n=== Project2 权限测试结果 ===")
	for _, tc := range testCases {
		ok, err := e.Enforce(tc.user, tc.domain, tc.path, tc.method)
		if err != nil {
			log.Printf("Enforce failed: %v", err)
			continue
		}
		fmt.Printf("用户: %-8s | 域: %-8s | 路径: %-25s | 方法: %-6s | 描述: %-30s | 允许访问: %v\n",
			tc.user, tc.domain, tc.path, tc.method, tc.desc, ok)
	}
}
