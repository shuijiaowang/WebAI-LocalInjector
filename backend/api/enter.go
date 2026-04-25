package api

import "Service/service"

// HandlerGroup 包含所有处理器的结构
type ApiGroup struct {
	ExampleApi
}

var (
	exampleService = service.ExampleService{}
)
