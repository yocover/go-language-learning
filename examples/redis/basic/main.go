package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "RedisP@ss2024!", // Redis 密码
		DB:       0,                // 使用默认数据库
	})

	// 创建上下文
	ctx := context.Background()

	// 1. String 操作
	fmt.Println("=== String 操作 ===")
	// 设置键值对
	err := rdb.Set(ctx, "name", "Alex", 0).Err()
	if err != nil {
		panic(err)
	}
	for i := range 10 {
		err = rdb.Set(ctx, fmt.Sprintf("name%d", i), fmt.Sprintf("Alex - value: %d", i), 0).Err()
		if err != nil {
			panic(err)
		}
	}

	// 获取值
	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("key = %v\n", val)

	// 设置过期时间
	err = rdb.Set(ctx, "temp_key", "temp_value", 15*time.Second).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("temp_key 将在 15 秒后过期")

	// 退出
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)

	// 2. Hash 操作
	fmt.Println("\n=== Hash 操作 ===")
	// 设置哈希表字段
	err = rdb.HSet(ctx, "user:1", "name", "John", "age", "25", "city", "New York").Err()
	if err != nil {
		panic(err)
	}

	// 获取所有字段
	fields, err := rdb.HGetAll(ctx, "user:1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("user:1 = %v\n", fields)

	// 3. List 操作
	fmt.Println("\n=== List 操作 ===")
	// 从左端推入元素
	err = rdb.LPush(ctx, "list", "first", "second", "third").Err()
	if err != nil {
		panic(err)
	}

	// 获取列表范围
	list, err := rdb.LRange(ctx, "list", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("list = %v\n", list)

	// 4. Set 操作
	fmt.Println("\n=== Set 操作 ===")
	// 添加集合元素
	err = rdb.SAdd(ctx, "set", "a", "b", "c", "a").Err()
	if err != nil {
		panic(err)
	}

	// 获取集合所有成员
	members, err := rdb.SMembers(ctx, "set").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("set = %v\n", members)

	// 5. Sorted Set 操作
	fmt.Println("\n=== Sorted Set 操作 ===")
	// 添加有序集合成员
	err = rdb.ZAdd(ctx, "scores", redis.Z{Score: 90, Member: "Alice"},
		redis.Z{Score: 85, Member: "Bob"},
		redis.Z{Score: 95, Member: "Carol"}).Err()
	if err != nil {
		panic(err)
	}

	// 获取有序集合范围（按分数从高到低）
	scores, err := rdb.ZRevRangeWithScores(ctx, "scores", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("scores (从高到低):")
	for _, z := range scores {
		fmt.Printf("  %v: %v\n", z.Member, z.Score)
	}

	// 6. 事务操作
	fmt.Println("\n=== 事务操作 ===")
	// 开始事务
	tx := rdb.TxPipeline()

	// 在事务中执行命令
	tx.Set(ctx, "tx_key1", "value1", 0)
	tx.Set(ctx, "tx_key2", "value2", 0)

	// 执行事务
	_, err = tx.Exec(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("事务执行成功")

	// 清理示例数据
	rdb.Del(ctx, "key", "temp_key", "user:1", "list", "set", "scores", "tx_key1", "tx_key2")
}
