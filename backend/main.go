package main

import (
	"Service/config"
	"Service/routes"
	"log"
	"strings"
)

func main() {

	config.InitConfig()
	r := routes.SetupRouter()

	port := config.AppConfig.Server.Port
	if strings.TrimSpace(port) == "" {
		port = "7789"
	}
	addr := ":" + strings.TrimPrefix(strings.TrimSpace(port), ":")
	// 启动服务
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
