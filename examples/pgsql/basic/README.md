# PostgreSQL Go示例

这个示例展示了如何在Go中使用PostgreSQL数据库，包括基本的CRUD操作和一些PostgreSQL特有的功能。

## 功能展示

1. 数据库连接
2. 表创建
3. 基本的CRUD操作（创建、读取、更新、删除）
4. 使用JSONB类型存储复杂数据
5. 高级查询示例（模糊搜索和JSON查询）

## 前置条件

1. 安装并运行PostgreSQL服务器
2. 创建测试数据库：
   ```sql
   CREATE DATABASE testdb;
   ```
3. 确保PostgreSQL服务器正在运行，并且可以使用以下默认配置连接：
   - 主机：localhost
   - 端口：5432
   - 用户名：postgres
   - 密码：postgres
   - 数据库：testdb

## 运行示例

1. 确保已安装Go和必要的依赖：
   ```bash
   go mod tidy
   ```

2. 运行示例：
   ```bash
   go run main.go
   ```

## 代码结构

- `User` 结构体：定义用户数据模型
- `PostgresDB` 结构体：封装数据库操作
- CRUD操作方法：
  - `CreateUser`：创建新用户
  - `GetUser`：获取用户信息
  - `UpdateUser`：更新用户信息
  - `DeleteUser`：删除用户
  - `SearchUsers`：搜索用户（支持模糊搜索和JSON查询）

## 注意事项

1. 示例中使用的连接字符串可能需要根据你的PostgreSQL配置进行修改
2. 示例使用JSONB类型存储用户配置文件，这是PostgreSQL的一个高级特性
3. 所有数据库操作都使用参数化查询，以防止SQL注入
4. 示例包含了错误处理，但在生产环境中可能需要更完善的错误处理机制 