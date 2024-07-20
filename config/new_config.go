package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig 加载并合并多个配置文件和环境变量到给定的结构体中
// paths中的配置文件按顺序读取，后面的配置文件会覆盖前面的配置文件
// 文件格式可以是json、yaml、toml
func LoadConfigFromFilesAndEnv(envPrefix string, files []string, config interface{}) error {
	v := viper.New()

	// 设置环境变量前缀
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取并合并配置文件
	for _, file := range files {
		v.SetConfigFile(file)
		if err := v.MergeInConfig(); err != nil {
			log.Printf("Error reading config file %s: %v", file, err)
		}
	}

	// 反序列化配置到结构体
	if err := v.Unmarshal(config); err != nil {
		return err
	}

	return nil
}
