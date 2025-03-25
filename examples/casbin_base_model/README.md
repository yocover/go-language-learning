# Casbin 基础权限模块

这是一个基于 Casbin 的通用权限管理模块，可以被其他项目复用。支持多种权限模型，使用 MySQL 存储权限规则。

## 特性

- 支持多种权限模型（RBAC、ABAC、ACL等）
- 支持多域隔离
- 使用 MySQL 存储权限规则
- 支持动态更新权限规则
- 线程安全
- 支持自动加载策略更新

## 安装

```bash
go get github.com/casbin/casbin/v2
go get github.com/casbin/gorm-adapter/v3
go get gorm.io/driver/mysql
go get gorm.io/gorm
```

## 数据库配置

1. 创建数据库：

```sql
CREATE DATABASE casbin_demo;
```

2. Casbin 会自动创建必要的表结构。

## 使用方法

1. 创建配置：

```go
config := &enforcer.Config{
    DBType:       "mysql",
    DBConnection: "user:password@tcp(127.0.0.1:3306)/casbin_demo?charset=utf8mb4&parseTime=True&loc=Local",
    ModelPath:    "path/to/your/model.conf",
    AutoLoad:     true,
}
```

2. 创建 enforcer：

```go
e, err := enforcer.NewEnforcer(config)
if err != nil {
    log.Fatalf("Failed to create enforcer: %v", err)
}
```

3. 添加策略：

```go
// 添加单个策略
e.AddPolicy("admin", "domain1", "/api/*", "*")

// 添加角色
e.AddGroupingPolicy("alice", "admin", "domain1")
```

4. 检查权限：

```go
ok, err := e.Enforce("alice", "domain1", "/api/users", "GET")
if err != nil {
    log.Printf("Check permission failed: %v", err)
}
fmt.Printf("Has permission: %v\n", ok)
```

## 示例项目

查看 `examples` 目录中的示例项目：

- `project1`: 简单的 API 权限控制示例
- `project2`: 部门级别的权限控制示例

## 自定义模型

你可以根据需要自定义权限模型。在 `models` 目录下提供了一些常用的模型配置：

- `rbac_with_domains.conf`: 支持多域的 RBAC 模型

## 注意事项

1. 确保正确配置数据库连接信息
2. 根据实际需求选择或自定义权限模型
3. 考虑是否需要启用自动加载功能
4. 在高并发环境下注意性能优化

## 性能优化建议

1. 合理设置自动加载间隔
2. 使用适当的缓存策略
3. 避免过于复杂的权限规则
4. 适当使用批量操作

## 常见问题

1. 数据库连接失败
   - 检查数据库连接字符串
   - 确保数据库服务正在运行
   - 验证用户名和密码

2. 权限规则不生效
   - 检查模型配置文件
   - 验证策略是否正确添加
   - 确认角色分配是否正确

3. 性能问题
   - 优化数据库索引
   - 调整自动加载间隔
   - 简化权限规则

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License 