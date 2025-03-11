package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// ExcelData 表示Excel解析后的数据结构
type ExcelData struct {
	RowID     int               `json:"row_id"`
	RowData   map[string]string `json:"row_data"`
	ParseTime string            `json:"parse_time"`
}

// ExcelDB 处理Excel数据的Redis操作
type ExcelDB struct {
	rdb       *redis.Client
	shardSize int    // 每个分片的大小
	prefix    string // Hash key的前缀
}

// NewExcelDB 创建ExcelDB实例
func NewExcelDB(rdb *redis.Client) *ExcelDB {
	return &ExcelDB{
		rdb:       rdb,
		shardSize: 1000, // 每个分片存储1000条数据
		prefix:    "excel_data",
	}
}

// getShardKey 获取分片key
func (e *ExcelDB) getShardKey(key string) string {
	// 使用key的前两位作为分片ID
	shardID := key[:2]
	return fmt.Sprintf("%s:%s", e.prefix, shardID)
}

// StoreData 存储Excel数据
func (e *ExcelDB) StoreData(ctx context.Context, key string, data []ExcelData) error {
	// 将数据转换为JSON字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal json error: %v", err)
	}

	// 获取分片key
	shardKey := e.getShardKey(key)

	// 使用 Hash 存储数据
	err = e.rdb.HSet(ctx, shardKey, key, jsonData).Err()
	if err != nil {
		return err
	}

	// 更新索引
	return e.rdb.SAdd(ctx, e.prefix+"_shards", shardKey).Err()
}

// GetData 获取Excel数据
func (e *ExcelDB) GetData(ctx context.Context, key string) ([]ExcelData, error) {
	// 获取分片key
	shardKey := e.getShardKey(key)

	// 从 Hash 获取数据
	jsonStr, err := e.rdb.HGet(ctx, shardKey, key).Result()
	if err != nil {
		return nil, err
	}

	// 解析JSON数据
	var data []ExcelData
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json error: %v", err)
	}

	return data, nil
}

// ListKeys 列出所有数据的键
func (e *ExcelDB) ListKeys(ctx context.Context) ([]string, error) {
	// 获取所有分片
	shards, err := e.rdb.SMembers(ctx, e.prefix+"_shards").Result()
	if err != nil {
		return nil, err
	}

	// 获取所有分片中的键
	var allKeys []string
	for _, shard := range shards {
		keys, err := e.rdb.HKeys(ctx, shard).Result()
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
	}

	return allKeys, nil
}

// DeleteData 删除指定的数据
func (e *ExcelDB) DeleteData(ctx context.Context, key string) error {
	// 获取分片key
	shardKey := e.getShardKey(key)

	// 删除数据
	return e.rdb.HDel(ctx, shardKey, key).Err()
}

// BatchStoreData 批量存储数据
func (e *ExcelDB) BatchStoreData(ctx context.Context, dataMap map[string][]ExcelData) error {
	// 按分片分组数据
	shardData := make(map[string]map[string][]byte)

	for key, data := range dataMap {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshal json error for key %s: %v", key, err)
		}

		shardKey := e.getShardKey(key)
		if _, ok := shardData[shardKey]; !ok {
			shardData[shardKey] = make(map[string][]byte)
		}
		shardData[shardKey][key] = jsonData
	}

	// 使用Pipeline批量存储每个分片的数据
	pipe := e.rdb.Pipeline()
	for shardKey, data := range shardData {
		for key, jsonData := range data {
			pipe.HSet(ctx, shardKey, key, jsonData)
		}
		pipe.SAdd(ctx, e.prefix+"_shards", shardKey)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// GetAllData 获取所有数据
func (e *ExcelDB) GetAllData(ctx context.Context) (map[string][]ExcelData, error) {
	// 获取所有分片
	shards, err := e.rdb.SMembers(ctx, e.prefix+"_shards").Result()
	if err != nil {
		return nil, err
	}

	// 获取所有分片中的数据
	dataMap := make(map[string][]ExcelData)
	for _, shard := range shards {
		result, err := e.rdb.HGetAll(ctx, shard).Result()
		if err != nil {
			return nil, err
		}

		for key, jsonStr := range result {
			var data []ExcelData
			err = json.Unmarshal([]byte(jsonStr), &data)
			if err != nil {
				return nil, fmt.Errorf("unmarshal json error for key %s: %v", key, err)
			}
			dataMap[key] = data
		}
	}

	return dataMap, nil
}

// GetShardStats 获取分片统计信息
func (e *ExcelDB) GetShardStats(ctx context.Context) (map[string]int64, error) {
	// 获取所有分片
	shards, err := e.rdb.SMembers(ctx, e.prefix+"_shards").Result()
	if err != nil {
		return nil, err
	}

	// 统计每个分片的数据量
	stats := make(map[string]int64)
	for _, shard := range shards {
		count, err := e.rdb.HLen(ctx, shard).Result()
		if err != nil {
			return nil, err
		}
		stats[shard] = count
	}

	return stats, nil
}

func main() {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "RedisP@ss2024!",
		DB:       0,
	})

	ctx := context.Background()
	excelDB := NewExcelDB(rdb)

	// 生成大量测试数据
	fmt.Println("生成测试数据...")
	batchData := make(map[string][]ExcelData)
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("%d::13", i)
		batchData[key] = []ExcelData{
			{
				RowID: i,
				RowData: map[string]string{
					"name": fmt.Sprintf("User%d", i),
					"age":  strconv.Itoa(20 + (i % 50)),
				},
				ParseTime: time.Now().Format("2006-01-02 15:04:05"),
			},
		}
	}

	// 批量存储数据
	fmt.Println("批量存储数据...")
	err := excelDB.BatchStoreData(ctx, batchData)
	if err != nil {
		log.Fatalf("批量存储失败: %v", err)
	}

	// 获取分片统计信息
	fmt.Println("\n分片统计信息:")
	stats, err := excelDB.GetShardStats(ctx)
	if err != nil {
		log.Fatalf("获取分片统计失败: %v", err)
	}
	for shard, count := range stats {
		fmt.Printf("分片 %s: %d 条数据\n", shard, count)
	}

	// 读取示例数据
	fmt.Println("\n读取示例数据:")
	key := "0::13"
	data, err := excelDB.GetData(ctx, key)
	if err != nil {
		log.Fatalf("读取数据失败: %v", err)
	}
	for _, row := range data {
		fmt.Printf("行ID: %d\n", row.RowID)
		fmt.Printf("数据: %v\n", row.RowData)
		fmt.Printf("解析时间: %s\n", row.ParseTime)
	}

	fmt.Println("\nPress Enter to exit...")
	os.Stdin.Read(make([]byte, 1))

	// 清理数据
	fmt.Println("\n清理数据...")
	shards, _ := excelDB.rdb.SMembers(ctx, excelDB.prefix+"_shards").Result()
	for _, shard := range shards {
		excelDB.rdb.Del(ctx, shard)
	}
	excelDB.rdb.Del(ctx, excelDB.prefix+"_shards")
	fmt.Println("数据已清理")
}
