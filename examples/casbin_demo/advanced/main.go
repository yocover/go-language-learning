package main

import (
	"fmt"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
)

/*
  # 管理员可以访问所有API和所有方法
   p, admin, /api/*, *

   # 开发者的权限
   p, developer, /api/v1/products/*, GET    # 查看所有产品
   p, developer, /api/v1/products/*, POST   # 创建产品
   p, developer, /api/v1/products/:id, PUT  # 更新特定产品
   p, developer, /api/v1/products/:id, DELETE # 删除特定产品

   # 普通用户的权限
   p, user, /api/v1/products/*, GET         # 只能查看产品
   p, user, /api/v1/products/:id/comments, POST # 可以发表评论

	 URL 匹配特性:
		支持精确匹配：/api/v1/products
		支持路径参数：/api/v1/products/:id
		支持通配符：/api/*
		支持嵌套资源：/api/v1/products/:id/comments
	权限级别:
		管理员（admin）：完全访问权限
		开发者（developer）：产品的完整 CRUD 操作
		普通用户（user）：只读权限和评论权限
*/

func ExampleAPIAccess() {
	fmt.Println("运行 API 访问控制示例...")
	e, err := casbin.NewEnforcer("api_model.conf", "api_policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer failed: %v", err)
	}

	// 测试用例
	testCases := []struct {
		user   string
		path   string
		method string
		desc   string
	}{
		// 管理员权限测试
		{"alice", "/api/v1/products", "GET", "管理员查看产品列表"},
		{"alice", "/api/v1/users", "POST", "管理员创建用户"},
		{"alice", "/api/v2/settings", "PUT", "管理员修改系统设置"},

		// 开发者权限测试
		{"bob", "/api/v1/products", "GET", "开发者查看产品列表"},
		{"bob", "/api/v1/products", "POST", "开发者创建新产品"},
		{"bob", "/api/v1/products/123", "PUT", "开发者更新特定产品"},
		{"bob", "/api/v1/products/123", "DELETE", "开发者删除特定产品"},
		{"bob", "/api/v2/settings", "GET", "开发者尝试访问系统设置"},

		// 普通用户权限测试
		{"charles", "/api/v1/products", "GET", "用户查看产品列表"},
		{"charles", "/api/v1/products/123", "GET", "用户查看特定产品"},
		{"charles", "/api/v1/products/123/comments", "POST", "用户发表评论"},
		{"charles", "/api/v1/products", "POST", "用户尝试创建产品"},
		{"charles", "/api/v1/products/123", "DELETE", "用户尝试删除产品"},
	}

	fmt.Println("\n=== 测试 REST API 访问控制 ===")
	fmt.Println("\n权限测试结果:")
	for _, tc := range testCases {
		ok, _ := e.Enforce(tc.user, tc.path, tc.method)
		fmt.Printf("用户: %-8s | 操作: %-6s | 路径: %-30s | %-20s | 允许访问: %v\n",
			tc.user,
			tc.method,
			tc.path,
			tc.desc,
			ok)
	}

	// 展示一些特殊的URL匹配案例
	fmt.Println("\n=== URL 模式匹配测试 ===")
	specialCases := []struct {
		user   string
		path   string
		method string
		desc   string
	}{
		{"bob", "/api/v1/products/999", "PUT", "动态参数URL"},
		{"bob", "/api/v1/products/special-item", "PUT", "带有特殊字符的URL"},
		{"charles", "/api/v1/products/123/comments", "POST", "嵌套资源URL"},
		{"charles", "/api/v1/products/456/comments/789", "GET", "深层嵌套URL"},
	}

	for _, tc := range specialCases {
		ok, _ := e.Enforce(tc.user, tc.path, tc.method)
		fmt.Printf("用户: %-8s | 操作: %-6s | 路径: %-30s | %-20s | 允许访问: %v\n",
			tc.user,
			tc.method,
			tc.path,
			tc.desc,
			ok)
	}
}

// 定义一个更丰富的用户结构体，用于ABAC示例
type User struct {
	Name       string
	Age        int
	IsAdmin    bool
	Department string
	Title      string
	JoinDate   time.Time
}

func (u *User) String() string {
	return u.Name
}

func ExampleHierarchicalRBAC() {
	e, err := casbin.NewEnforcer("hierarchical_rbac_model.conf", "hierarchical_rbac_policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer failed: %v", err)
	}

	fmt.Println("\n=== 测试多层角色继承 ===")
	// 测试不同层级的权限
	checkPermission := func(sub, obj, act string) {
		ok, _ := e.Enforce(sub, obj, act)
		fmt.Printf("%s %s data 的权限 %s: %v\n", sub, act, obj, ok)
	}

	fmt.Println("\nAlice (admin) 的权限:")
	checkPermission("alice", "data", "write")
	checkPermission("alice", "data", "read")
	checkPermission("alice", "data", "view")

	fmt.Println("\nBob (manager) 的权限:")
	checkPermission("bob", "data", "write")
	checkPermission("bob", "data", "read")
	checkPermission("bob", "data", "view")

	fmt.Println("\nCharles (user) 的权限:")
	checkPermission("charles", "data", "write")
	checkPermission("charles", "data", "read")
	checkPermission("charles", "data", "view")
}

func ExampleResourceHierarchy() {
	e, err := casbin.NewEnforcer("resource_hierarchy_model.conf", "resource_hierarchy_policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer failed: %v", err)
	}

	fmt.Println("\n=== 测试资源层级 ===")
	checkAccess := func(sub, obj, act string) {
		ok, _ := e.Enforce(sub, obj, act)
		fmt.Printf("%s 对 %s 的 %s 权限: %v\n", sub, obj, act, ok)
	}

	fmt.Println("\n测试 Alice (admin) 的权限:")
	checkAccess("alice", "/data/project/secret", "write")
	checkAccess("alice", "/data/public/docs", "write")

	fmt.Println("\n测试 Bob (developer) 的权限:")
	checkAccess("bob", "/data/project/secret", "read")
	checkAccess("bob", "/data/project/public", "read")
	checkAccess("bob", "/data/public/docs", "read")

	fmt.Println("\n测试 Charles (user) 的权限:")
	checkAccess("charles", "/data/public/docs", "read")
	checkAccess("charles", "/data/project/secret", "read")
}

func ExampleDomainIsolation() {
	e, err := casbin.NewEnforcer("domain_model.conf", "domain_policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer failed: %v", err)
	}

	fmt.Println("\n=== 测试域/租户隔离 ===")
	checkDomainAccess := func(sub, dom, obj, act string) {
		ok, _ := e.Enforce(sub, dom, obj, act)
		fmt.Printf("%s 在域 %s 中对 %s 的 %s 权限: %v\n", sub, dom, obj, act, ok)
	}

	fmt.Println("\n测试跨域权限:")
	checkDomainAccess("alice", "domain1", "data", "write")
	checkDomainAccess("alice", "domain2", "data", "write")
	checkDomainAccess("catherine", "domain1", "data", "write")
	checkDomainAccess("catherine", "domain2", "data", "write")
}

func ExampleABAC() {
	e, err := casbin.NewEnforcer("abac_model.conf", "abac_policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer failed: %v", err)
	}

	fmt.Println("\n=== 测试ABAC（基于属性的访问控制）===")

	// 创建不同属性的用户
	alice := &User{
		Name:       "alice",
		Age:        35,
		IsAdmin:    true,
		Department: "IT",
		Title:      "Senior Engineer",
		JoinDate:   time.Now().AddDate(-5, 0, 0), // 5年前加入
	}

	bob := &User{
		Name:       "bob",
		Age:        20,
		IsAdmin:    false,
		Department: "IT",
		Title:      "Junior Engineer",
		JoinDate:   time.Now().AddDate(0, -6, 0), // 6个月前加入
	}

	charlie := &User{
		Name:       "charlie",
		Age:        40,
		IsAdmin:    false,
		Department: "HR",
		Title:      "Manager",
		JoinDate:   time.Now().AddDate(-2, 0, 0), // 2年前加入
	}

	// 添加各种基于属性的策略
	rules := []struct {
		sub string
		obj string
		act string
	}{
		// 规则1: IT部门的高级工程师可以访问代码库
		{"r.sub.Department == 'IT' && r.sub.Title == 'Senior Engineer'", "code_repository", "write"},

		// 规则2: 任何IT部门的成员可以读取文档
		{"r.sub.Department == 'IT'", "documents", "read"},

		// 规则3: 管理层可以访问财务报告
		{"r.sub.Title == 'Manager'", "financial_reports", "read"},

		// 规则4: 工作超过2年的员工可以访问内部系统
		{"time.Now().Sub(r.sub.JoinDate).Hours() > 17520", "internal_systems", "access"},

		// 规则5: 综合条件：25岁以上的IT部门成员或任何管理层可以参加决策会议
		{"(r.sub.Age > 25 && r.sub.Department == 'IT') || r.sub.Title == 'Manager'", "decision_meeting", "attend"},
	}

	// 添加所有规则
	for _, rule := range rules {
		_, err = e.AddPolicy(rule.sub, rule.obj, rule.act)
		if err != nil {
			log.Printf("AddPolicy failed: %v", err)
		}
	}

	// 测试不同用户的权限
	testCases := []struct {
		user *User
		obj  string
		act  string
		desc string
	}{
		{alice, "code_repository", "write", "高级工程师访问代码库"},
		{bob, "code_repository", "write", "初级工程师访问代码库"},
		{alice, "documents", "read", "IT成员访问文档"},
		{charlie, "documents", "read", "非IT成员访问文档"},
		{charlie, "financial_reports", "read", "管理层访问财务报告"},
		{bob, "financial_reports", "read", "非管理层访问财务报告"},
		{alice, "internal_systems", "access", "老员工访问内部系统"},
		{bob, "internal_systems", "access", "新员工访问内部系统"},
		{alice, "decision_meeting", "attend", "高级IT成员参加决策会议"},
		{charlie, "decision_meeting", "attend", "管理层参加决策会议"},
		{bob, "decision_meeting", "attend", "年轻IT成员参加决策会议"},
	}

	fmt.Println("\n权限测试结果:")
	for _, tc := range testCases {
		ok, _ := e.Enforce(tc.user, tc.obj, tc.act)
		fmt.Printf("%s (%s, %s) %s: %v\n",
			tc.user.Name,
			tc.user.Title,
			tc.user.Department,
			tc.desc,
			ok)
	}
}

func main() {
	fmt.Println("运行高级 Casbin 示例...")
	ExampleHierarchicalRBAC()
	ExampleResourceHierarchy()
	ExampleDomainIsolation()
	ExampleABAC()
	ExampleAPIAccess()
}
