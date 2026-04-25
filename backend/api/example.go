package api

import (
	"Service/model/request"
	"Service/util/response"

	"github.com/gin-gonic/gin"
)

type ExampleApi struct{}

func (h *ExampleApi) Test(c *gin.Context) {
	var req request.ExampleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Result(7, "ceshi", "成功", c)
		return
	}
	response.Result(0, "ceshi——失败", "失败", c)
}
