package controller

import (
	"fmt"
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

	rootUrl := configs.Service.ModulesDir + env.Separator + name + env.Separator + version
	for i, file := range files {
		filepath := filepaths[i]
		f, err1 := file.Open()
		if err1 != nil {
			ctx.ErrorParams(err1.Error())
			return
		}
		defer f.Close()

		targetFilepath := rootUrl + env.Separator + filepath
		err2 := utils.FileCreateWithDirs(targetFilepath)
		if err2 != nil {
			ctx.ErrorCustom(err2.Error())
			return
		}

		data, err2 := io.ReadAll(f)
		if err2 != nil {
			ctx.ErrorParams(err2.Error())
			return
		}
		err3 := os.WriteFile(targetFilepath, data, 0644)
		if err3 != nil {
			ctx.ErrorParams(err3.Error())
			return
		}

		fmt.Println(file, filepath)
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
		ctx.ErrorCustom("package not exist")
		return
	}

	err := os.RemoveAll(rootUrl)
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}

	// 空目录
	files, _ := os.ReadDir(packageFilePath)
	if len(files) == 0 {
		os.RemoveAll(packageFilePath)
	}

	ctx.Success()
}
