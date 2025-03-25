package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

/*
这个示例展示了基本的RBAC模型，但Casbin还支持更复杂的场景，比如：
多层角色继承
资源层级
域/租户隔离
基于属性的访问控制（ABAC）
等等
*/

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		log.Printf("Enforce failed: %v", err)
		return
	}
	if ok {
		fmt.Printf("%s 可以 %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s 不能 %s %s\n", sub, act, obj)
	}
}

func main() {
	// 加载模型和策略
	e, err := casbin.NewEnforcer("rbac_model.conf", "rbac_policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer failed: %v", err)
	}

	// 测试权限
	fmt.Println("测试 alice 的权限 (admin角色):")
	check(e, "alice", "data1", "read")
	check(e, "alice", "data1", "write")
	check(e, "alice", "data2", "read")
	check(e, "alice", "data2", "write")

	fmt.Println("\n测试 bob 的权限 (user角色):")
	check(e, "bob", "data1", "read")
	check(e, "bob", "data1", "write")
	check(e, "bob", "data2", "read")
	check(e, "bob", "data2", "write")

	// 添加新的角色和权限
	fmt.Println("\n动态添加权限:")
	_, err = e.AddPolicy("user", "data2", "read")
	if err != nil {
		log.Printf("AddPolicy failed: %v", err)
	}
	fmt.Println("为user角色添加了data2的读权限")

	fmt.Println("\n再次测试 bob 的权限:")
	check(e, "bob", "data2", "read")
}
