package main

import (
	"fmt"
	"gdp/command/handler"
	"os"
	"strings"
)

func main() {
	cmd := ""
	args := os.Args
	for i := 1; i < len(args); i++ {
		cmd += args[i] + " "
	}

	// 安装
	if strings.HasPrefix(cmd, "install") {
		packageName, version := getPackageInfo(cmd)
		handler.Install(packageName, version)
		return
	}

	// 卸载
	if strings.HasPrefix(cmd, "uninstall") {
		packageName, _ := getPackageInfo(cmd)
		handler.Uninstall(packageName)
		return
	}

	// 发布
	if strings.HasPrefix(cmd, "publish") {
		handler.Publish()
		return
	}

	// 取消发布
	if strings.HasPrefix(cmd, "unpublish") {
		packageName, version := getPackageInfo(cmd)
		fmt.Println(packageName, version)
		return
	}

	fmt.Println("unknown command")
}

func getPackageInfo(cmd string) (string, string) {
	split1 := strings.Split(cmd, " ")
	packageName := split1[1]
	split2 := strings.Split(packageName, "@")
	version := ""
	if len(split2) == 2 {
		packageName = split2[0]
		version = split2[1]
	}
	return packageName, version
}
