package controller

import (
	"gdp/service/configs"
	"gdp/service/middleware"
	"gdp/utils"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// 是否允许该路径
func fileAllowed(filename string) bool {
	for _, v := range configs.Service.AllowedCatalogs {
		if strings.HasPrefix(filename, v) {
			return true
		}
	}
	return false
}

// 获取文件目录
func FileCatalog(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	type Params struct {
		Path string `form:"path" binding:"required"`
		Deep bool   `form:"deep"`
	}
	params := Params{}
	if err := c.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}
	if !fileAllowed(params.Path) {
		ctx.ErrorAuth("You don't have permission to access this file")
		return
	}

	list, err := utils.FileCatalog(params.Path, params.Deep)
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}
	ctx.SuccessData(list)
}

func FileInfo(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	type Params struct {
		Path string `form:"path" binding:"required"`
		Deep bool   `form:"deep"`
	}
	params := Params{}
	if err := c.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}
	if !fileAllowed(params.Path) {
		ctx.ErrorAuth("You don't have permission to access this file")
		return
	}

	info, err := utils.FileInfo(params.Path, params.Deep)
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}
	ctx.SuccessData(info)
}

// 获取文件内容
func FileContent(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	type Params struct {
		Path string `form:"path" binding:"required"`
	}
	params := Params{}
	if err := c.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}
	if !fileAllowed(params.Path) {
		ctx.ErrorAuth("You don't have permission to access this file")
		return
	}

	content, err := os.ReadFile(params.Path)
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}
	ctx.Data(200, "text/plain", content)
	ctx.Abort()
}
