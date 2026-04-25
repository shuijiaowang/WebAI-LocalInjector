package api

import (
	"Service/model/request"
	"Service/service"
	"Service/util/response"

	"github.com/gin-gonic/gin"
)

type FileTreeApi struct{}

var fileTreeService = service.FileTreeService{}

func (h *FileTreeApi) GetFileTree(c *gin.Context) {
	var req request.FileTreeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Result(400, nil, "参数错误: "+err.Error(), c)
		return
	}

	tree, err := fileTreeService.GetFileTree(req.RootPath, req.IgnoreDirs, req.IgnoreFiles, req.IgnoreExts)
	if err != nil {
		response.Result(500, nil, "读取目录失败: "+err.Error(), c)
		return
	}

	response.Result(0, tree, "成功", c)
}
