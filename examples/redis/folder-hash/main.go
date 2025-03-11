package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// FileMetadata 文件元数据结构
type FileMetadata struct {
	Bucket       string `json:"bucket"`
	PathPrefix   string `json:"path_prefix"`
	ID           string `json:"id"`
	GeoWorkspace string `json:"geo_workspace"`
	ObjectKey    string `json:"object_key"`
	FileName     string `json:"file_name"`
}

// FolderDB 模拟文件夹结构的 Redis 操作
type FolderDB struct {
	rdb *redis.Client
}

// NewFolderDB 创建文件夹数据库实例
func NewFolderDB(rdb *redis.Client) *FolderDB {
	return &FolderDB{rdb: rdb}
}

// CreateFile 创建文件及其元数据
func (f *FolderDB) CreateFile(ctx context.Context, folderPath string, fileID string, metadata FileMetadata) error {
	// 创建完整的键名
	key := fmt.Sprintf("%s::%s", folderPath, fileID)

	// 将结构体转换为 map[string]interface{}
	metadataMap := make(map[string]interface{})
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &metadataMap)
	if err != nil {
		return err
	}

	// 使用 HSet 存储所有字段
	return f.rdb.HSet(ctx, key, metadataMap).Err()
}

// GetFile 获取文件元数据
func (f *FolderDB) GetFile(ctx context.Context, folderPath string, fileID string) (*FileMetadata, error) {
	key := fmt.Sprintf("%s::%s", folderPath, fileID)

	// 获取所有字段
	result, err := f.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, redis.Nil
	}

	// 将 map 转换回结构体
	metadata := &FileMetadata{}
	jsonStr, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonStr, metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// ListFolder 列出文件夹内容
func (f *FolderDB) ListFolder(ctx context.Context, folderPath string) ([]string, error) {
	pattern := fmt.Sprintf("%s::*", folderPath)
	var keys []string
	var cursor uint64
	for {
		var err error
		var batch []string
		batch, cursor, err = f.rdb.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

// DeleteFile 删除文件
func (f *FolderDB) DeleteFile(ctx context.Context, folderPath string, fileID string) error {
	key := fmt.Sprintf("%s::%s", folderPath, fileID)
	return f.rdb.Del(ctx, key).Err()
}

func main() {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "RedisP@ss2024!",
		DB:       0,
	})

	ctx := context.Background()
	folderDB := NewFolderDB(rdb)

	// 创建示例文件
	metadata := FileMetadata{
		Bucket:       "lyd-base",
		PathPrefix:   "路径前缀",
		ID:           "179749028347576320",
		GeoWorkspace: "lyd-base/neuearth/open-data/179749028347576320",
		ObjectKey:    "neuearth/open-data/179749028347576320",
		FileName:     "The northern area of Jeddah City.tif",
	}

	// 创建文件结构
	folderPath := "compliance_analysis_topic/excel/18621110343895777729/[Empty]"
	fileID := "13"

	fmt.Println("创建文件...")
	err := folderDB.CreateFile(ctx, folderPath, fileID, metadata)
	if err != nil {
		panic(err)
	}

	// 获取文件信息
	fmt.Println("\n读取文件信息:")
	file, err := folderDB.GetFile(ctx, folderPath, fileID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("文件名: %s\n", file.FileName)
	fmt.Printf("Bucket: %s\n", file.Bucket)
	fmt.Printf("Object Key: %s\n", file.ObjectKey)

	// 列出文件夹内容
	fmt.Println("\n列出文件夹内容:")
	files, err := folderDB.ListFolder(ctx, folderPath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Printf("- %s\n", f)
	}

	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// os.Exit(0)

	// 删除文件
	fmt.Println("\n删除文件...")
	err = folderDB.DeleteFile(ctx, folderPath, fileID)
	if err != nil {
		panic(err)
	}
	fmt.Println("文件已删除")

}
