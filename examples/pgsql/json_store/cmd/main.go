package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"example/pgsql/json_store/store"
)

func main() {
	// 数据库连接配置
	connStr := "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable"
	fmt.Printf("Connecting to database with connection string: %s\n", connStr)

	// 创建store实例
	st, err := store.NewStore(connStr)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer st.Close()

	fmt.Println("Successfully connected to database")

	ctx := context.Background()

	// 创建表
	fmt.Println("Creating table if not exists...")
	if err := st.CreateTables(ctx); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// 清理旧数据（如果存在）
	if err := st.DeleteFileJSON(ctx, "test-file-1"); err != nil {
		// 忽略不存在的记录错误
		if !os.IsNotExist(err) {
			log.Printf("Warning: Failed to clean up old data: %v", err)
		}
	}

	// 示例JSON数据
	largeJSON := map[string]interface{}{
		"id": "123",
		"metadata": map[string]interface{}{
			"name":      "example",
			"type":      "test",
			"timestamp": time.Now().Unix(),
		},
		"data": map[string]interface{}{
			"field1": "value1",
			"field2": 42,
			"field3": []string{"a", "b", "c"},
			"nested": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	// 将JSON数据转换为字符串
	jsonStr, err := json.MarshalIndent(largeJSON, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// 创建文件JSON记录
	fileJSON := &store.FileJSON{
		FileID:      "test-file-1",
		FileName:    "test.json",
		JSONContent: string(jsonStr),
	}

	// 存储JSON
	if err := st.CreateFileJSON(ctx, fileJSON); err != nil {
		log.Fatalf("Failed to create file JSON: %v", err)
	}
	fmt.Printf("Created file JSON with ID: %d\n", fileJSON.ID)

	// 读取JSON
	retrieved, err := st.GetFileJSON(ctx, fileJSON.FileID)
	if err != nil {
		log.Fatalf("Failed to get file JSON: %v", err)
	}

	// 解析JSON内容
	var parsedJSON map[string]interface{}
	if err := json.Unmarshal([]byte(retrieved.JSONContent), &parsedJSON); err != nil {
		log.Fatalf("Failed to parse JSON content: %v", err)
	}

	fmt.Printf("\nRetrieved and parsed JSON for file %s:\n", retrieved.FileName)
	prettyJSON, _ := json.MarshalIndent(parsedJSON, "", "  ")
	fmt.Println(string(prettyJSON))

	// 更新JSON内容
	largeJSON["metadata"].(map[string]interface{})["updated"] = true
	updatedJSONStr, _ := json.MarshalIndent(largeJSON, "", "  ")

	fileJSON.JSONContent = string(updatedJSONStr)
	if err := st.UpdateFileJSON(ctx, fileJSON); err != nil {
		log.Fatalf("Failed to update file JSON: %v", err)
	}
	fmt.Printf("\nUpdated file JSON at: %v\n", fileJSON.UpdatedAt)

	// 列出所有文件（不包含JSON内容）
	files, err := st.ListFiles(ctx)
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	fmt.Printf("\nAll files in storage:\n")
	for _, f := range files {
		fmt.Printf("- %s (ID: %s, Created: %v)\n", f.FileName, f.FileID, f.CreatedAt)
	}

	// 等待用户输入后再删除数据
	fmt.Println("\nPress Enter to clean up and exit...")
	os.Stdin.Read(make([]byte, 1))

	// 删除文件
	if err := st.DeleteFileJSON(ctx, fileJSON.FileID); err != nil {
		log.Printf("Warning: Failed to delete file JSON: %v", err)
	} else {
		fmt.Printf("Deleted file with ID: %s\n", fileJSON.FileID)
	}
}
