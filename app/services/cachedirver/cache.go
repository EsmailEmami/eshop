package cachedirver

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/spf13/viper"
)

var conn CacheDriver

func GetConnection() CacheDriver {
	if conn == nil {
		Use("redis")
	}

	return conn
}

func Use(driverName string) {
	switch driverName {
	case "redis":
		{
			address := viper.GetString("redis.host") + ":" + viper.GetString("redis.port")
			db := viper.GetInt("redis.db")
			redisConn := redis.NewClient(&redis.Options{
				Addr:     address,
				Password: viper.GetString("redis.password"),
				DB:       db,
			})

			pool := goredis.NewPool(redisConn)
			rs := redsync.New(pool)
			mutex := rs.NewMutex("redismutex")

			conn = &RedisDriver{mutex: mutex}
			conn.SetConnection(redisConn)

			return
		}
	}
	panic("invalid cache driver")
}
