package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

// Leaderboard Redis 排行榜实现
type Leaderboard struct {
	rdb *redis.Client
	key string
}

// NewLeaderboard 创建排行榜
func NewLeaderboard(rdb *redis.Client, key string) *Leaderboard {
	return &Leaderboard{
		rdb: rdb,
		key: fmt.Sprintf("leaderboard:%s", key),
	}
}

// AddScore 添加或更新分数
func (l *Leaderboard) AddScore(ctx context.Context, member string, score float64) error {
	return l.rdb.ZAdd(ctx, l.key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// GetRank 获取排名（从0开始）
func (l *Leaderboard) GetRank(ctx context.Context, member string) (int64, error) {
	// 使用 ZRevRank 获取倒序排名（分数高的排前面）
	return l.rdb.ZRevRank(ctx, l.key, member).Result()
}

// GetScore 获取分数
func (l *Leaderboard) GetScore(ctx context.Context, member string) (float64, error) {
	return l.rdb.ZScore(ctx, l.key, member).Result()
}

// GetTopN 获取前 N 名
func (l *Leaderboard) GetTopN(ctx context.Context, n int64) ([]redis.Z, error) {
	return l.rdb.ZRevRangeWithScores(ctx, l.key, 0, n-1).Result()
}

// GetRankRange 获取指定排名范围的成员
func (l *Leaderboard) GetRankRange(ctx context.Context, start, stop int64) ([]redis.Z, error) {
	return l.rdb.ZRevRangeWithScores(ctx, l.key, start, stop).Result()
}

func main() {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "RedisP@ss2024!",
		DB:       0,
	})

	ctx := context.Background()

	// 创建游戏排行榜
	leaderboard := NewLeaderboard(rdb, "game")

	// 模拟玩家得分
	players := []string{"player1", "player2", "player3", "player4", "player5"}
	fmt.Println("添加玩家分数...")

	rand.Seed(time.Now().UnixNano())
	for _, player := range players {
		// 随机生成分数 (0-1000)
		score := float64(rand.Intn(1000))
		err := leaderboard.AddScore(ctx, player, score)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s 的分数: %.0f\n", player, score)
	}

	// 获取前 3 名
	fmt.Println("\n排行榜前三名:")
	top3, err := leaderboard.GetTopN(ctx, 3)
	if err != nil {
		panic(err)
	}

	for i, z := range top3 {
		fmt.Printf("第 %d 名: %s (分数: %.0f)\n", i+1, z.Member, z.Score)
	}

	// 获取特定玩家的排名
	player := "player1"
	rank, err := leaderboard.GetRank(ctx, player)
	if err != nil {
		panic(err)
	}
	score, err := leaderboard.GetScore(ctx, player)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n玩家 %s 的排名: 第 %d 名 (分数: %.0f)\n", player, rank+1, score)

	// 获取排名范围
	fmt.Println("\n所有玩家排名:")
	allRanks, err := leaderboard.GetRankRange(ctx, 0, -1)
	if err != nil {
		panic(err)
	}

	for i, z := range allRanks {
		fmt.Printf("第 %d 名: %s (分数: %.0f)\n", i+1, z.Member, z.Score)
	}

	// 更新玩家分数
	fmt.Println("\n更新玩家分数...")
	player = "player1"
	newScore := float64(rand.Intn(1000))
	err = leaderboard.AddScore(ctx, player, newScore)
	if err != nil {
		panic(err)
	}
	fmt.Printf("更新 %s 的分数为: %.0f\n", player, newScore)

	// 获取更新后的排名
	newRank, err := leaderboard.GetRank(ctx, player)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s 的新排名: 第 %d 名\n", player, newRank+1)

	// 清理排行榜数据
	err = rdb.Del(ctx, leaderboard.key).Err()
	if err != nil {
		panic(err)
	}
}
