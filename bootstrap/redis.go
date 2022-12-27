package bootstrap

import (
	"charites/global"
	"context"

	redis "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func setupRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     global.RedisSetting.Address,
		Password: global.RedisSetting.Password, // 密码
		DB:       global.RedisSetting.DB,       // 数据库
		PoolSize: global.RedisSetting.PoolSize, // 连接池大小
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	pool := goredis.NewPool(client)
	global.Redsync = redsync.New(pool)

	return nil
}
