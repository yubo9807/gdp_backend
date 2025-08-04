package configs

import (
	"gdp/env"
	"os"

	"gopkg.in/yaml.v2"
)

type ServiceType struct {
	Port            int
	Prefix          string
	LogDir          string   `yaml:"logDir"`
	LogReserveTime  int      `yaml:"logReserveTime"`
	ModulesDir      string   `yaml:"modulesDir"`
	AllowedCatalogs []string `yaml:"allowedCatalogs"`
}

var Service ServiceType

const serviceTemplate = `
prefix: "/api"
port: 9800  # 启动端口

logDir: "./logs"  # 日志目录
logReserveTime: 10  # 日志保留时间(d)

modulesDir: "./modules"  # 依赖包目录
allowedCatalogs: ["modules"]  # 允许接口访问的目录
`

func init() {
	configFile := "./" + env.ServiceConfigName
	data, err := os.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
		os.WriteFile(configFile, []byte(serviceTemplate), 0777)
		data, _ = os.ReadFile(configFile)
	}

	if err := yaml.Unmarshal([]byte(data), &Service); err != nil {
		panic(err.Error())
	}
}
