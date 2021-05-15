package zhenai

import (
	"fmt"
	"test/model"

	"github.com/go-redis/redis"
)

//爬虫去重 ，添加的是有序集合
func PutCityData(url string) error {

	// 添加有序集合 插入成功为1 插入失败为0,去重用
	value, err := model.Pool.ZAdd("zhenai_city_spider", redis.Z{Score: 10, Member: url}).Result()
	fmt.Println("zhenai,城市url添加", url, err)

	if value == 1 { //说明没有这个key
		onlyid := model.WxHead + model.GetKey(16)
		//存对应的data到有序集合
		value, err := model.Pool.ZAdd("zhenai_get_data", redis.Z{Score: 10, Member: onlyid}).Result()
		fmt.Println(value, err)
		ma := map[string]interface{}{}
		ma["url"] = url
		//再存入map参数
		err = model.Pool.HMSet(onlyid, ma).Err()
		if err != nil {
			fmt.Println("PutCityData", err)
		}
	}
	return err
}

func GetCityData() (url_list []string, err error) {

	//设置最大和最小值  返回有序集合的所有元素和分数
	vals, err := model.Pool.ZRangeByScoreWithScores("zhenai_get_data", redis.ZRangeBy{
		Min:    "0",
		Max:    "50",
		Offset: 0,
		Count:  1,
	}).Result()

	for _, value := range vals {
		key := value.Member.(string)
		dMap, err := model.Pool.HGetAll(key).Result()
		if err != nil {
			return url_list, err
		}
		url_list = append(url_list, dMap["url"])
		model.Pool.Del(key).Err()
		//删除集合中的一个指定元素
		model.Pool.ZRem("zhenai_get_data", key)
	}

	return url_list, err
}
