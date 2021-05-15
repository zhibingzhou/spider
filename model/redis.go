package model

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

type RedisConf struct {
	Host   string
	Port   string
	Pwd    string
	DBName int
}

var Head = "all_key"
var WxHead = "all_key_wx"
var Pool *redis.Client
var ctx = context.Background()

// redis初始化
func InitRedis(redisMsg RedisConf) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     redisMsg.Host + ":" + redisMsg.Port,
		Password: redisMsg.Pwd,
		DB:       redisMsg.DBName,
	})
	err := client.Ping().Err()
	if err != nil {
		log.Fatalln(err)
	}
	Pool = client
	return client
}

func init() {
	Pool = InitRedis(RedisConf{
		Host:   "127.0.0.1",
		Port:   "6379",
		Pwd:    "foobared",
		DBName: 0,
	})
}

//清理缓存
func Delcash() error {

	//拿到key头在集合中的数量
	num, err := Pool.SCard(Head).Result()
	if err != nil {
		return err
	}
	var i int64
	for i = 0; i < num; i++ {

		//删除一条数据返回被删除的元素，逐个删除，但这个会返回对应元素
		red_key, err := Pool.SPop(Head).Result()

		if err != nil {
			return err
		}

		if Pool.Del(red_key).Err() != nil {
			return err
		}

	}

	return err
}
