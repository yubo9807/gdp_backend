package controller

import (
	"archive/zip"
	"gdp/env"
	"gdp/service/configs"
	"gdp/service/middleware"
	"gdp/utils"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 获取包的最大版本
func getDirectoryMaxVersion(filename string) (string, error) {
	files, err := os.ReadDir(filename)
	if err != nil {
		return "", err
	}
	fileNames := []string{}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return utils.VersionMax(fileNames), nil
}

func Install(c *gin.Context) {
	ctx := middleware.ContextGet(c)

	type Params struct {
		Name    string `form:"name" binding:"required"`
		Version string `form:"version"`
	}
	params := Params{}
	if err := c.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	packageFilePath := configs.Service.ModulesDir + env.Separator + params.Name
	if params.Version == "" {
		maxVersion, _ := getDirectoryMaxVersion(packageFilePath)
		params.Version = maxVersion
	} else {
		_, err := os.Stat(packageFilePath + env.Separator + params.Version)
		if err != nil {
			ctx.ErrorCustom("version not exist")
			return
		}
	}
	rootUrl := packageFilePath + env.Separator + params.Version

	size := int64(0)
	paths := []string{}
	filepath.Walk(rootUrl, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		size += info.Size()
		paths = append(paths, path)
		return nil
	})

	ctx.Header("Content-Length", strconv.FormatInt(size, 10))
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+params.Name+".zip")
	ctx.Header("Content-Type", "application/octet-stream")

	zipWriter := zip.NewWriter(ctx.Writer)
	defer zipWriter.Close()

	for _, path := range paths {
		file, _ := os.Open(path)
		defer file.Close()
		url := params.Name + strings.Replace(path, rootUrl[2:], "", -1)
		zipFile, err := zipWriter.Create(url)
		if err != nil {
			break
		}
		io.Copy(zipFile, file)
	}

}
