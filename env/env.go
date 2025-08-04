package env

import (
	"os"
	"path/filepath"
)

const (
	ProjectName = "gdp"
	Separator   = string(filepath.Separator)

	PackageConfigName = ProjectName + ".json"
	ServiceConfigName = "config_service.yml"
	CommandConfigName = "config_command.yml"
	UserDataFileName  = "./db/user.json"
)

var CommandDir string

func init() {
	home, _ := os.UserHomeDir()
	CommandDir = home + Separator + ProjectName
}
