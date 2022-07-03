package cmn

import (
	"github.com/95eh/eg"
	"github.com/gomodule/redigo/redis"
	"time"
)

func GetRedisFac(conf RedisConf) (connFac eg.ToRedisConnError, connPool *redis.Pool) {
	redisFac := func() (redis.Conn, error) {
		return redis.Dial("tcp", conf.Addr,
			redis.DialPassword(conf.Password),
			redis.DialDatabase(conf.Db))
	}
	redisPool := &redis.Pool{
		Dial:        redisFac,
		IdleTimeout: 300 * time.Second,
		MaxActive:   512,
		MaxIdle:     512,
	}
	return redisFac, redisPool
}
