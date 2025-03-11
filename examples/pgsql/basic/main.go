package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// User 表示用户数据结构
type User struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Profile   json.RawMessage `json:"profile"`
	CreatedAt time.Time       `json:"created_at"`
}

// PostgresDB 处理数据库操作
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresDB 创建PostgresDB实例
func NewPostgresDB(connStr string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// 测试连接
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

// CreateTables 创建必要的表
func (p *PostgresDB) CreateTables(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		profile JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := p.db.ExecContext(ctx, query)
	return err
}

// CreateUser 创建新用户
func (p *PostgresDB) CreateUser(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (name, email, profile)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	return p.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Profile).
		Scan(&user.ID, &user.CreatedAt)
}

// GetUser 获取用户信息
func (p *PostgresDB) GetUser(ctx context.Context, id int) (*User, error) {
	user := &User{}
	query := `
	SELECT id, name, email, profile, created_at
	FROM users
	WHERE id = $1`

	err := p.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Profile, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (p *PostgresDB) UpdateUser(ctx context.Context, user *User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, profile = $3
	WHERE id = $4`

	result, err := p.db.ExecContext(ctx, query, user.Name, user.Email, user.Profile, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}
	return nil
}

// DeleteUser 删除用户
func (p *PostgresDB) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

// SearchUsers 搜索用户（演示高级查询）
func (p *PostgresDB) SearchUsers(ctx context.Context, searchTerm string) ([]User, error) {
	query := `
	SELECT id, name, email, profile, created_at
	FROM users
	WHERE 
		name ILIKE $1 OR 
		email ILIKE $1 OR 
		profile @> $2::jsonb`

	// 创建JSON查询条件
	jsonQuery := fmt.Sprintf(`{"interests": ["%s"]}`, searchTerm)

	rows, err := p.db.QueryContext(ctx, query, "%"+searchTerm+"%", jsonQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Profile, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func main() {
	// 连接字符串
	connStr := "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable"

	// 创建数据库连接
	db, err := NewPostgresDB(connStr)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	defer db.db.Close()

	ctx := context.Background()

	// 创建表
	if err := db.CreateTables(ctx); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	// 创建用户示例
	profile := json.RawMessage(`{
		"interests": ["coding", "reading"],
		"location": "Beijing",
		"age": 25
	}`)

	user := &User{
		Name:    "张三",
		Email:   "zhangsan@example.com",
		Profile: profile,
	}

	// 创建用户
	err = db.CreateUser(ctx, user)
	if err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}
	fmt.Printf("创建的用户ID: %d\n", user.ID)

	// 获取用户
	fetchedUser, err := db.GetUser(ctx, user.ID)
	if err != nil {
		log.Fatalf("获取用户失败: %v", err)
	}
	fmt.Printf("获取的用户信息: %+v\n", fetchedUser)

	// 更新用户
	user.Name = "李四"
	err = db.UpdateUser(ctx, user)
	if err != nil {
		log.Fatalf("更新用户失败: %v", err)
	}
	fmt.Println("用户更新成功")

	// 搜索用户
	users, err := db.SearchUsers(ctx, "coding")
	if err != nil {
		log.Fatalf("搜索用户失败: %v", err)
	}
	fmt.Printf("找到 %d 个用户喜欢编程\n", len(users))

	// 删除用户
	err = db.DeleteUser(ctx, user.ID)
	if err != nil {
		log.Fatalf("删除用户失败: %v", err)
	}
	fmt.Println("用户删除成功")
}
