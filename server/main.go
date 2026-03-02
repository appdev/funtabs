package main

import (
	"flag"
	"fmt"
	"funtabs-server/config"
	"funtabs-server/model"
	"funtabs-server/server"
	"funtabs-server/storage"
	"log"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 1. 加载配置
	config.Load(*cfgPath)

	// 2. 初始化数据库
	model.Init()

	// 3. 初始化存储后端
	storage.Init()

	// 4. 启动 HTTP 服务
	r := server.NewRouter(config.Cfg.Server.Origins)
	addr := fmt.Sprintf(":%s", config.Cfg.Server.Port)
	log.Printf("funtabs-server 启动在 %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
