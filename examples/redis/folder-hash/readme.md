从截图中我看到的是一个 Redis Desktop Manager 的界面，展示了一个不同的文件夹结构组织方式。这种结构使用了 Redis Hash 数据类型来存储数据，具体特点是：

1. **层级结构**:
   ```
   compliance_analysis_topic/
   └── excel/
       └── 18621110343895777729/
           └── [Empty]/
               └── excel:18621110343895777729::13
   ```

2. **数据存储格式**:
   - 使用 Hash 类型存储键值对
   - 键包含多个字段：
     - `bucket`: "lyd-base"
     - `路径前缀`: (中文字段)
     - `id`: "179749028347576320"
     - `geo_workspace`: "lyd-base/neuearth/open-data/..."
     - `object_key`: "neuearth/open-data/..."
     - `file_name`: "The northern area of Jeddah City.tif"

让我修改之前的代码来实现这种结构：

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

// FileMetadata 文件元数据结构
type FileMetadata struct {
	Bucket        string `json:"bucket"`
	PathPrefix    string `json:"path_prefix"`
	ID            string `json:"id"`
	GeoWorkspace  string `json:"geo_workspace"`
	ObjectKey     string `json:"object_key"`
	FileName      string `json:"file_name"`
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

	// 删除文件
	fmt.Println("\n删除文件...")
	err = folderDB.DeleteFile(ctx, folderPath, fileID)
	if err != nil {
		panic(err)
	}
	fmt.Println("文件已删除")
}
```
这个新版本的代码实现了与截图中相似的文件夹结构，主要特点包括：

1. **文件元数据结构**:
   ```go
   type FileMetadata struct {
       Bucket        string `json:"bucket"`
       PathPrefix    string `json:"path_prefix"`
       ID            string `json:"id"`
       GeoWorkspace  string `json:"geo_workspace"`
       ObjectKey     string `json:"object_key"`
       FileName      string `json:"file_name"`
   }
   ```

2. **键名格式**:
   - 使用双冒号（`::`）作为分隔符
   - 格式：`folderPath::fileID`
   - 例如：`compliance_analysis_topic/excel/18621110343895777729/[Empty]::13`

3. **主要功能**:
   - `CreateFile`: 创建文件及其元数据，使用 Redis Hash 存储
   - `GetFile`: 获取文件的完整元数据
   - `ListFolder`: 列出指定文件夹下的所有文件
   - `DeleteFile`: 删除指定文件

4. **数据存储方式**:
   - 使用 Redis Hash 类型存储每个文件的元数据
   - 每个字段都可以单独获取和更新
   - 支持完整的元数据检索

要运行这个示例：

```bash
cd examples/redis/folder
go run main.go
```

输出将展示：
1. 创建一个示例文件
2. 读取文件的元数据信息
3. 列出文件夹内容
4. 删除文件

这种实现方式的优点：
1. 更好的数据结构化
2. 支持复杂的元数据
3. 高效的字段访问
4. 符合实际业务需求

需要注意的是，这个实现使用了 Redis Hash 类型，这样可以更有效地存储和检索文件的元数据信息，与截图中显示的结构保持一致。
