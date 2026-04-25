package routes

import (
	"Service/api"
	"Service/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 跨域中间件（放在最前面）
	r.Use(middleware.Cors())
	exampleApi := api.ExampleApi{}
	fileTreeApi := api.FileTreeApi{}
	// 用户路由
	// 消费记录路由（需要认证）
	apiGroup := r.Group("/api")
	exampleGroup := apiGroup.Group("/example")
	{
		exampleGroup.POST("/test", exampleApi.Test) // 添加消费记录
	}
	fileGroup := apiGroup.Group("/file")
	{
		fileGroup.POST("/tree", fileTreeApi.GetFileTree)
	}
	return r
}
