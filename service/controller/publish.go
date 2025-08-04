package controller

import (
	"encoding/json"
	"gdp/env"
	"gdp/service/configs"
	"gdp/service/middleware"
	"gdp/utils"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadFile struct {
	Filename    string    `json:"filename"`
	FilePath    string    `json:"filepath"`
	UUID        string    `json:"uuid"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType"`
	UploadTime  time.Time `json:"uploadTime"`
}

// 发布
func Publish(c *gin.Context) {
	ctx := middleware.ContextGet(c)

	form, err := c.MultipartForm()
	if err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	name := form.Value["name"][0]
	version := form.Value["version"][0]
	files := form.File["file"]
	filepaths := form.Value["filepath"]

	// 当前账号有无权限
	packageName := configs.Service.ModulesDir + env.Separator + name
	powerFilename := packageName + env.Separator + env.PackagePowerName
	userInfo := GetUserInfo(c)
	if _, err := os.Stat(powerFilename); err != nil {
		// 第一次发布版本
		utils.FileCreateWithDirs(powerFilename)
		powerConfig := PowerConfig{
			Author:        userInfo.Name,
			Collaborators: []string{},
		}
		content, _ := json.Marshal(powerConfig)
		os.WriteFile(powerFilename, content, 0644)
	} else {
		// 后期维护，检查协作者
		content, _ := os.ReadFile(powerFilename)
		powerConfig := PowerConfig{}
		json.Unmarshal(content, &powerConfig)
		powerConfig.Collaborators = append(powerConfig.Collaborators, userInfo.Name)
		isCollaborators := utils.SliceSome(powerConfig.Collaborators, func(v string, i int) bool {
			return v == userInfo.Name
		})
		if !isCollaborators {
			ctx.ErrorAuth("你还不是这个包的协作者")
			return
		}
	}

	rootUrl := configs.Service.ModulesDir + env.Separator + name + env.Separator + version
	for i, file := range files {
		filepath := filepaths[i]
		f, openErr := file.Open()
		if openErr != nil {
			ctx.ErrorParams(openErr.Error())
			return
		}
		defer f.Close()

		targetFilepath := rootUrl + env.Separator + filepath
		utils.FileCreateWithDirs(targetFilepath)

		data, readErr := io.ReadAll(f)
		if readErr != nil {
			ctx.ErrorParams(readErr.Error())
			return
		}
		writeErr := os.WriteFile(targetFilepath, data, 0644)
		if writeErr != nil {
			ctx.ErrorParams(writeErr.Error())
			return
		}
	}

	ctx.Success()
}

// 取消发布
func Unpublish(c *gin.Context) {
	ctx := middleware.ContextGet(c)

	type Params struct {
		Name    string `json:"name" binding:"required"`
		Version string `json:"version" binding:"required"`
	}
	params := Params{}
	if err := c.ShouldBindJSON(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	packageFilePath := configs.Service.ModulesDir + env.Separator + params.Name
	rootUrl := packageFilePath + env.Separator + params.Version
	_, statErr := os.Stat(rootUrl)
	if statErr != nil {
		ctx.ErrorCustom(params.Name + ": " + params.Version + " 版本已存在")
		return
	}

	err := os.RemoveAll(rootUrl)
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}

	// 空目录
	files, _ := os.ReadDir(packageFilePath)
	if len(files) == 1 {
		os.RemoveAll(packageFilePath)
	}

	ctx.Success()
}
