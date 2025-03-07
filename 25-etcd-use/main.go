package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// 创建etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 设置上下文
	ctx := context.Background()

	// 1.用户信息，使用层级结构
	cli.Put(ctx, "users/1/info", `
	{"name": "张三", "age": 25, "email": "zhangsan@example.com"}
	`)
	cli.Put(ctx, "users/1/settings", `
	{"theme": "dark", "language": "zh-CN"}
	`)

	// 2.用户状态信息，使用状态前缀
	cli.Put(ctx, "status/users/1", "online")
	cli.Put(ctx, "status/users/2", "offline")

	// 3.用户权限信息，使用权限前缀
	cli.Put(ctx, "permissions/users/1", `["read", "write", "admin"]`)

	// 4.用户会话信息，使用租约（临时数据）
	lease, _ := cli.Grant(ctx, 10)                                       // 租约有效期10秒
	cli.Put(ctx, "sessions/1", "token123", clientv3.WithLease(lease.ID)) // 写入租约ID

	// 5. 用户索引，方便查找
	cli.Put(ctx, "indexes/email/zhangsan@example.com", "1")
	cli.Put(ctx, "indexes/phone/13800138000", "1")

	// 查询示例：获取某个用户的所有信息
	resp, _ := cli.Get(ctx, "users/1/", clientv3.WithPrefix())
	for _, ev := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", ev.Key, ev.Value)
	}

}

func testUser(cli *clientv3.Client, ctx context.Context) {

	var err error

	// 写入简直对
	_, err = cli.Put(ctx, "name", "张三")
	if err != nil {
		log.Fatal(err)
	}

	// 读取值
	resp, err := cli.Get(ctx, "name")
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", ev.Key, ev.Value)
	}

	// 查询所有用户
	resp, err = cli.Get(ctx, "user/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", ev.Key, ev.Value)
	}

	watchChan := cli.Watch(ctx, "name")
	go func() {
		for wresp := range watchChan {
			for _, ev := range wresp.Events {
				fmt.Printf("Type: %s Key: %s Value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	// 删除键
	_, err = cli.Delete(ctx, "name")
	if err != nil {
		log.Fatal(err)
	}

	// 等待一会儿以观察变化
	time.Sleep(time.Second)

	// 写入多个带前缀的键值对
	cli.Put(ctx, "user/1", "张三")
	cli.Put(ctx, "user/2", "李四")
	cli.Put(ctx, "user/3", "王五")
}
