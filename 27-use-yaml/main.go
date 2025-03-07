package main

import (
	"example/config/config"
	"fmt"
	"log"
)

func main() {

	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// 使用配置
	fmt.Printf("服务器地址： %s:%d\n", cfg.Server.Host, cfg.Server.Port)

	fmt.Printf("数据库地址： %s@%s:%d/%s\n",
		cfg.Database.Username,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)
}
