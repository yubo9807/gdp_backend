package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gdp/command/configs"
	"gdp/env"
	"gdp/utils"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Name    string
	Version string
	Author  string
	Files   []string
	Modules map[string]string
}

func Publish() {
	packageFilename := "./" + env.PackageConfigName
	content, err := os.ReadFile(packageFilename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	config := Config{}
	err2 := json.Unmarshal(content, &config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}

	// 遍历文件
	paths := []string{packageFilename}
	isAll := len(config.Files) == 0
	filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isAll {
			paths = append(paths, path)
		} else {
			bl := utils.SliceSome(config.Files, func(v string, i int) bool {
				return v == path || strings.HasPrefix(path, v+env.Separator)
			})
			if bl {
				paths = append(paths, path)
			}
		}
		return nil
	})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", config.Name)
	writer.WriteField("version", config.Version)
	for _, url := range paths {
		file, err := os.Open(url)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		io.Copy(part, file)
		writer.WriteField("filepath", url)
	}
	writer.Close()
	fmt.Println(body)

	res, err := utils.Request("POST", configs.Command.RequestUrl+"/publish", map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if res.StatusCode == 200 {
		fmt.Println("Publish success")
	} else {
		fmt.Println(res)
	}
}
