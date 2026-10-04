package core

import (
	"log"
	"server/config"
	"server/utils"

	"gopkg.in/yaml.v3"
)

// 把磁盘上的`config.yaml` 文件解析成内存中的`Config` 结构体·
func InitConf() *config.Config {
	// 初始化配置结构体
	c := &config.Config{}

	// 读取配置文件
	yamlConfig, err := utils.LoadYAML()
	if err != nil {
		log.Fatal("Failed to load config.yaml", err)
		return nil
	}

	//`yaml.Unmarshal` 根据各字段的`yaml:"..."` 标签（如`yaml:"mysql"` ）把 YAML 里每个小节填进`Mysql` 、`Redis` 、`Jwt` 等子结构体字段
	err = yaml.Unmarshal(yamlConfig, c)
	if err != nil {
		log.Fatal("Failed to unmarshal config.yaml", err)
		return nil
	}
	return c
}
