package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// User 用户结构体
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	LastSeen time.Time `json:"last_seen"`
}

// 模拟数据库操作
func getUserFromDB(id int) (*User, error) {
	// 这里模拟从数据库查询数据
	// 实际应用中，这里会查询真实的数据库
	time.Sleep(200 * time.Millisecond) // 模拟数据库查询延迟
	return &User{
		ID:       id,
		Name:     "John Doe",
		Email:    "john@example.com",
		LastSeen: time.Now(),
	}, nil
}

func main() {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "RedisP@ss2024!",
		DB:       0,
	})

	ctx := context.Background()

	// 用户 ID
	userID := 1
	cacheKey := fmt.Sprintf("user:%d", userID)

	// 1. 尝试从缓存获取用户数据
	fmt.Println("尝试从缓存获取用户数据...")
	userJSON, err := rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		// 缓存未命中，从数据库获取
		fmt.Println("缓存未命中，从数据库获取...")
		user, err := getUserFromDB(userID)
		if err != nil {
			panic(err)
		}

		// 将用户数据序列化为 JSON
		userJSON, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		// 将数据存入缓存，设置过期时间为 1 小时
		err = rdb.Set(ctx, cacheKey, userJSON, time.Hour).Err()
		if err != nil {
			panic(err)
		}
		fmt.Println("用户数据已缓存")

		// 打印用户数据
		fmt.Printf("用户数据 (从数据库): %+v\n", user)
	} else if err != nil {
		panic(err)
	} else {
		// 缓存命中，解析数据
		var user User
		err := json.Unmarshal([]byte(userJSON), &user)
		if err != nil {
			panic(err)
		}
		fmt.Println("缓存命中！")
		fmt.Printf("用户数据 (从缓存): %+v\n", user)
	}

	// 2. 演示缓存更新
	fmt.Println("\n更新用户数据...")
	updatedUser := &User{
		ID:       userID,
		Name:     "John Doe Updated",
		Email:    "john.updated@example.com",
		LastSeen: time.Now(),
	}

	// 序列化更新后的用户数据
	updatedJSON, err := json.Marshal(updatedUser)
	if err != nil {
		panic(err)
	}

	// 更新缓存
	err = rdb.Set(ctx, cacheKey, updatedJSON, time.Hour).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("缓存已更新")

	// 3. 验证更新
	fmt.Println("\n验证更新后的数据...")
	newUserJSON, err := rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		panic(err)
	}

	var newUser User
	err = json.Unmarshal([]byte(newUserJSON), &newUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("更新后的用户数据: %+v\n", newUser)

	// 4. 清理缓存
	fmt.Println("\n清理缓存...")
	err = rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("缓存已清理")
}
