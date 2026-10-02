package utils

import (
	"io/fs"
	"os"
	"server/global"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

// 给initConfig()提供原始字节
func LoadYAML() ([]byte, error) {
	return os.ReadFile(configFile)
}

func SaveYAML() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, byteData, fs.ModePerm)
}
