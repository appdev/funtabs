package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	DB      DBConfig      `mapstructure:"db"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Storage StorageConfig `mapstructure:"storage"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Origins string `mapstructure:"origins"`
}

type DBConfig struct {
	Path string `mapstructure:"path"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // 单位：天
}

type StorageConfig struct {
	Type      string   `mapstructure:"type"`       // "local" 或 "s3"
	LocalPath string   `mapstructure:"local_path"` // 本地存储目录
	CDN       string   `mapstructure:"cdn"`        // CDN 前缀
	S3        S3Config `mapstructure:"s3"`
}

type S3Config struct {
	Endpoint  string `mapstructure:"endpoint"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// Cfg 是全局配置实例
var Cfg Config

// Load 从指定路径加载配置文件，支持环境变量覆盖
func Load(path string) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败 (%s): %v", path, err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}
}
