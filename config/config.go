package config

import (
	"fmt"
	"go-template/pkg/filepath"
	"os"
	"time"

	"github.com/spf13/viper"
)

var Conf = viper.New()

func init() {
	defaultConfig()
	homePath := os.Getenv("HOME")
	filePath := fmt.Sprintf( // 配置文件默认为 $HOME/.config/${project_name}/settings.toml
		"%s/.config/%s/settings.toml",
		homePath, Conf.GetString("project_name"),
	)
	Conf.SetConfigFile(filePath)
	if !filepath.IsExist(filePath) {
		filepath.CreateFile(filePath)
		Conf.WriteConfig()
	}
	// 如果还是不存在则结束程序, 并打印错误
	if err := Conf.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found, %s", err))
		} else {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}
	// 监听配置文件变化
	Conf.WatchConfig()
}

func defaultConfig() {
	Conf.SetDefault("project_name", "go-template") // 改成自己的名字吧
	// 默认设置, 推荐在配置文件中修改
	// 服务器设置
	Conf.SetDefault("host.address", "127.0.0.1")
	Conf.SetDefault("host.port", "8080")
	Conf.SetDefault("host.mode", "debug")
	Conf.SetDefault("host.tokenExpireDuration", 72*time.Hour)
	// redis
	Conf.SetDefault("redis.host", "localhost:6379")
	Conf.SetDefault("redis.password", "")
	Conf.SetDefault("redis.dbnum", 0)
	Conf.SetDefault("redis.poolSize", 10)
	Conf.SetDefault("redis.minIdleConns", 5)
	Conf.SetDefault("redis.dialTimeout", 5*time.Second)
	Conf.SetDefault("redis.readTimeout", 3*time.Second)
	Conf.SetDefault("redis.writeTimeout", 4*time.Second)
	Conf.SetDefault("redis.PoolTimeout", 4*time.Second)
	Conf.SetDefault("redis.idleCheckFrequency", 60*time.Second)
	Conf.SetDefault("redis.idleTimeout", 5*time.Minute)
	Conf.SetDefault("redis.maxConnAge", 0*time.Second)
	Conf.SetDefault("redis.maxRetries", 2)
	Conf.SetDefault("redis.minRetryBackoff", 8*time.Millisecond)
	Conf.SetDefault("redis.maxRetryBackoff", 512*time.Millisecond)
	// mysql
	Conf.SetDefault("mysql.host", "localhost")
	Conf.SetDefault("mysql.port", "3306")
	Conf.SetDefault("mysql.username", "root")
	Conf.SetDefault("mysql.password", "")
	Conf.SetDefault("mysql.dbName", "demo")
	Conf.SetDefault("mysql.maxIdleConns", 20)
	Conf.SetDefault("mysql.maxOpenConns", 20)
	Conf.SetDefault("mysql.connMaxLifetime", 3600)
}
