package env

import (
	"os"
	"path/filepath"
)

const (
	ProjectName = "gdp"
	Separator   = string(filepath.Separator)

	PackageConfigName = ProjectName + ".json"
	PackagePowerName  = "power.json"
	ServiceConfigName = "config_service.yml"
	CommandConfigName = "config_command.yml"
)

var CommandDir string

func init() {
	home, _ := os.UserHomeDir()
	CommandDir = home + Separator + ProjectName
}
