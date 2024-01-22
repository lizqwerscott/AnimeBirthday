package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v3/log"
)

type DatabaseConfig struct {
	Url      string `toml:"url"`
	Database string `toml:"database"`
	Coll     string `toml:"coll"`
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	Source   string `toml:"source"`
}

type ServerConfig struct {
	Port int `toml:"port"`
}

type Config struct {
	DBConfig     DatabaseConfig `toml:"db_config"`
	ServerConfig ServerConfig   `toml:"server_config"`
}

// LoadConfig 函数根据环境加载配置
func LoadConfig() Config {
	var configFilePath string
	env := os.Getenv("APP_ENV")

	switch env {
	case "development":
		configFilePath = "config.dev.toml"
	case "production":
		configFilePath = "config.prod.toml"
	default:
		configFilePath = "config.toml"
	}

	config, err := loadConfigFromFile(configFilePath)
	if err != nil {
		log.Errorf("加载配置时出错: %s", err)
		// 默认配置
		config = Config{
			DBConfig: DatabaseConfig{
				Url:      "mongodb://localhost:27017",
				Database: "AnimeBirthday",
				Coll:     "animebirthdays",
				User:     "animebirthday",
				Pass:     "123456",
				Source:   "AnimeBirthday",
			},
			ServerConfig: ServerConfig{
				Port: 22400,
			},
		}
	}

	return config
}

// loadConfigFromFile 函数从 TOML 文件中加载配置
func loadConfigFromFile(filePath string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
