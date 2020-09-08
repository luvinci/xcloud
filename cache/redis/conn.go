package redis

import (
	"time"
	cfg "xcloud/config"

	"github.com/gomodule/redigo/redis"
)

var (
	host        = cfg.Viper.GetString("redis.host")
	pwd         = cfg.Viper.GetString("redis.pwd")
	maxIdle     = cfg.Viper.GetInt("redis.maxIdle")
	maxActive   = cfg.Viper.GetInt("redis.maxActive")
	idleTimeout = cfg.Viper.GetDuration("redis.idleTimeout")
)

var pool *redis.Pool

func Pool() *redis.Pool {
	return pool
}

func init() {
	pool = newRedisPool()
}

// newRedisPool: 创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			// 1.打开连接
			conn, err = redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			// 2.访问认证
			if _, err = conn.Do("auth", pwd); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		// TestOnBorrow: 用于在应用程序再次使用该连接之前检查空闲连接的运行状况。 参数t是连接返回到池的时间。 如果函数返回错误，则连接将关闭。
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("ping")
			return err
		},
	}
}
