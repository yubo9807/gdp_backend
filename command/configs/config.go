package configs

import (
	"gdp/env"
	"gdp/utils"
	"os"

	"gopkg.in/yaml.v2"
)

type CommandType struct {
	PackageDir     string `yaml:"packageDir"`
	RequestUrl     string `yaml:"requestUrl"`
	LogReserveTime int    `yaml:"logReserveTime"`
}

const commandTemplate = `
packageDir: "./modules"
RequestUrl: "http://gdp.hpyyb.cn/api"
logReserveTime: 7
`

// 获取命令配置文件路径
func GetCommandConfigFilePath() string {
	filename := "./" + env.CommandConfigName
	_, err := os.Stat(filename)
	if err != nil {
		return env.CommandDir + env.Separator + env.CommandConfigName
	}
	return filename
}

var Command CommandType

func init() {
	configFile := GetCommandConfigFilePath()
	data, err := os.ReadFile(configFile)
	if err != nil {
		utils.FileCreateWithDirs(configFile)
		os.WriteFile(configFile, []byte(commandTemplate), 0777)
		data, _ = os.ReadFile(configFile)
	}

	if err := yaml.Unmarshal([]byte(data), &Command); err != nil {
		panic(err.Error())
	}
}
