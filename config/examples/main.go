package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig 加载并合并多个配置文件和环境变量到给定的结构体中
// paths中的配置文件按顺序读取，后面的配置文件会覆盖前面的配置文件
// 文件格式可以是json、yaml、toml
func LoadConfig(paths []string, config interface{}) error {
	v := viper.New()

	// 设置环境变量前缀
	v.SetEnvPrefix("MYAPP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取并合并配置文件
	for _, path := range paths {
		v.SetConfigFile(path)
		if err := v.MergeInConfig(); err != nil {
			log.Printf("Error reading config file %s: %v", path, err)
		}
	}

	// 反序列化配置到结构体
	if err := v.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

// Config 是一个表示应用程序配置的结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig 表示服务器配置
type ServerConfig struct {
	Port int
}

// DatabaseConfig 表示数据库配置
type DatabaseConfig struct {
	User     string
	Password string
	Name     string
}

func main() {
	configPaths := []string{
		"config1.yaml",
		"config2.toml",
	}

	var cfg Config
	if err := LoadConfig(configPaths, &cfg); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("Database User: %s\n", cfg.Database.User)
	fmt.Printf("Database Password: %s\n", cfg.Database.Password)
	fmt.Printf("Database Name: %s\n", cfg.Database.Name)
}
