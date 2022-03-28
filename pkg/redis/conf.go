package redis

import (
	"go-template/config"

	"github.com/go-redis/redis/v8"
)

var redisConfig = &redis.Options{
	Network:  "tcp",
	Addr:     config.Conf.GetString("redis.host"),
	Password: config.Conf.GetString("redis.password"),
	DB:       config.Conf.GetInt("redis.dbnum"),

	PoolSize:     config.Conf.GetInt("redis.poolSize"),
	MinIdleConns: config.Conf.GetInt("redis.minIdleConns"),

	DialTimeout:  config.Conf.GetDuration("redis.dialTimeout"),
	ReadTimeout:  config.Conf.GetDuration("redis.readTimeout"),
	WriteTimeout: config.Conf.GetDuration("redis.writeTimeout"),
	PoolTimeout:  config.Conf.GetDuration("redis.PoolTimeout"),

	IdleCheckFrequency: config.Conf.GetDuration("redis.idleCheckFrequency"),
	IdleTimeout:        config.Conf.GetDuration("redis.idleTimeout"),
	MaxConnAge:         config.Conf.GetDuration("redis.maxConnAge"),

	MaxRetries:      config.Conf.GetInt("redis.maxRetries"),
	MinRetryBackoff: config.Conf.GetDuration("redis.minRetryBackoff"),
	MaxRetryBackoff: config.Conf.GetDuration("redis.maxRetryBackoff"),
}
