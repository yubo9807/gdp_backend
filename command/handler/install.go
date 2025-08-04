package handler

import (
	"fmt"
	"gdp/command/configs"
	"gdp/env"
	"gdp/utils"
	"io"
	"os"
)

type WriteCounter struct {
	Len   int64
	Total int64
}

const format = "\rdownloading... %d \t total: %d"

func (w *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	w.Total += int64(n)
	fmt.Printf(format, w.Total, w.Len)
	return n, nil
}

func Install(packageName, version string) {
	packageDir := configs.Command.PackageDir
	os.RemoveAll(packageDir + env.Separator + packageName)

	requestUrl := configs.Command.RequestUrl + "/install" + "?name=" + packageName + "&version=" + version
	res, err := utils.Request("GET", requestUrl, nil, nil)
	if err != nil {
		fmt.Println("install error: ", err)
		return
	}

	outputFile, err := os.Create(packageDir + env.Separator + packageName + ".zip")
	if err != nil {
		fmt.Println("install error: ", err)
		return
	}

	counter := &WriteCounter{}
	counter.Len = res.ContentLength
	newa := io.TeeReader(res.Body, counter)
	io.Copy(outputFile, newa)
	outputFile.Close()

	utils.FileUnzip(outputFile.Name(), packageDir)
	os.Remove(outputFile.Name())

	fmt.Printf(format, counter.Len, counter.Len)
	fmt.Println("\ninstall success")
}

// 移除安装包
func Uninstall(packageName string) {
	filename := configs.Command.PackageDir + env.Separator + packageName
	os.RemoveAll(filename)
	fmt.Println("uninstall success")
}
