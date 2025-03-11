package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// FileJSON 表示文件JSON存储的数据结构
type FileJSON struct {
	ID          int       `json:"id"`
	FileID      string    `json:"file_id"`
	FileName    string    `json:"file_name"`
	JSONContent string    `json:"json_content"` // 使用string类型存储JSON内容
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Store 处理JSON数据存储的数据库操作
type Store struct {
	db *sql.DB
}

// NewStore 创建Store实例
func NewStore(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// CreateTables 创建必要的表
func (s *Store) CreateTables(ctx context.Context) error {
	// 先检查表是否存在
	checkQuery := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'dm_excel_json_store'
		);`

	var exists bool
	err := s.db.QueryRowContext(ctx, checkQuery).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check table existence failed: %v", err)
	}

	if exists {
		fmt.Println("Table dm_excel_json_store already exists")
		return nil
	}

	query := `
	CREATE TABLE IF NOT EXISTS dm_excel_json_store (
		id SERIAL PRIMARY KEY,
		file_id VARCHAR(255) NOT NULL UNIQUE,
		file_name VARCHAR(255) NOT NULL,
		json_content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("create table failed: %v", err)
	}

	fmt.Println("Table dm_excel_json_store created successfully")
	return nil
}

// CreateFileJSON 创建新的文件JSON记录
func (s *Store) CreateFileJSON(ctx context.Context, fileJSON *FileJSON) error {
	query := `
		INSERT INTO dm_excel_json_store (file_id, file_name, json_content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return s.db.QueryRowContext(ctx, query,
		fileJSON.FileID,
		fileJSON.FileName,
		fileJSON.JSONContent,
	).Scan(&fileJSON.ID, &fileJSON.CreatedAt, &fileJSON.UpdatedAt)
}

// GetFileJSON 通过file_id获取文件JSON记录
func (s *Store) GetFileJSON(ctx context.Context, fileID string) (*FileJSON, error) {
	fileJSON := &FileJSON{}
	query := `
		SELECT id, file_id, file_name, json_content, created_at, updated_at
		FROM dm_excel_json_store
		WHERE file_id = $1`

	err := s.db.QueryRowContext(ctx, query, fileID).
		Scan(
			&fileJSON.ID,
			&fileJSON.FileID,
			&fileJSON.FileName,
			&fileJSON.JSONContent,
			&fileJSON.CreatedAt,
			&fileJSON.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}
	return fileJSON, nil
}

// UpdateFileJSON 更新文件JSON记录
func (s *Store) UpdateFileJSON(ctx context.Context, fileJSON *FileJSON) error {
	query := `
		UPDATE dm_excel_json_store
		SET file_name = $1, json_content = $2, updated_at = CURRENT_TIMESTAMP
		WHERE file_id = $3
		RETURNING updated_at`

	return s.db.QueryRowContext(ctx, query,
		fileJSON.FileName,
		fileJSON.JSONContent,
		fileJSON.FileID,
	).Scan(&fileJSON.UpdatedAt)
}

// DeleteFileJSON 删除文件JSON记录
func (s *Store) DeleteFileJSON(ctx context.Context, fileID string) error {
	query := `DELETE FROM dm_excel_json_store WHERE file_id = $1`

	result, err := s.db.ExecContext(ctx, query, fileID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("file with id %s not found", fileID)
	}
	return nil
}

// ListFiles 列出所有文件基本信息（不包含JSON内容）
func (s *Store) ListFiles(ctx context.Context) ([]FileJSON, error) {
	query := `
		SELECT id, file_id, file_name, created_at, updated_at
		FROM dm_excel_json_store
		ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []FileJSON
	for rows.Next() {
		var fj FileJSON
		err := rows.Scan(
			&fj.ID,
			&fj.FileID,
			&fj.FileName,
			&fj.CreatedAt,
			&fj.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, fj)
	}
	return results, rows.Err()
}

// Close 关闭数据库连接
func (s *Store) Close() error {
	return s.db.Close()
}
