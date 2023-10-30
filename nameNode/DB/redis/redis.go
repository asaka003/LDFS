package redis

import (
	"LDFS/nameNode/config"
	"fmt"

	"github.com/go-redis/redis"
)

/*
	redis五大数据类型

	字符串(string)
	列表(List)
	集合(Set)
	哈希(Hash)
	有序集合(Zset)
*/

var RDB *redis.Client

func Close() {
	_ = RDB.Close()
}

func RedisInit() error {
	redis_opention := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,       //选择哪一个reids库，默认为0
		PoolSize: config.Redis.PoolSize, //连接池的大小
	}
	RDB = redis.NewClient(&redis_opention)
	_, err := RDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
