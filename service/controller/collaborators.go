package controller

import (
	"encoding/json"
	"gdp/env"
	"gdp/service/configs"
	"gdp/service/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

type PowerConfig struct {
	Author        string   `json:"author"`
	Collaborators []string `json:"collaborators"`
}

func Collaborators(c *gin.Context) {
	ctx := middleware.ContextGet(c)

	type Params struct {
		PackageName   string   `json:"packageName" binding:"required"`
		Collaborators []string `json:"collaborators" binding:"required"`
	}
	var params Params
	if err := c.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	filename := configs.Service.ModulesDir + env.Separator + params.PackageName + env.Separator + env.PackagePowerName
	content, err := os.ReadFile(filename)
	if err != nil {
		ctx.ErrorCustom("找不到发布包")
		return
	}
	powerConfig := PowerConfig{}
	jsonErr := json.Unmarshal(content, &powerConfig)
	if jsonErr != nil {
		ctx.ErrorCustom(jsonErr.Error())
		return
	}
	adminInfo := GetUserInfo(c)
	if adminInfo.Name != powerConfig.Author {
		ctx.ErrorCustom("您没有权限添加协作者")
		return
	}
	powerConfig.Collaborators = params.Collaborators
	newContent, _ := json.Marshal(powerConfig)
	os.WriteFile(filename, newContent, 0644)

	ctx.Success()
}
