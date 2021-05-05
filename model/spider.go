package model

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//爬虫去重 ，添加的是有序集合
func PutCityData(url string) error {

	// 添加有序集合 插入成功为1 插入失败为0,去重用
	value, err := pool.ZAdd("city_spider", redis.Z{Score: 10, Member: url}).Result()
	fmt.Println("城市url添加", url, err)

	if value == 1 { //说明没有这个key
		onlyid := wxhead + GetKey(16)
		//存对应的data到有序集合
		value, err := pool.ZAdd("get_data", redis.Z{Score: 10, Member: onlyid}).Result()
		fmt.Println(value, err)
		ma := map[string]interface{}{}
		ma["url"] = url
		//再存入map参数
		err = pool.HMSet(onlyid, ma).Err()
		if err != nil {
			fmt.Println("PutCityData", err)
		}
	}
	return err
}

func GetCityData() (url_list []string, err error) {

	//设置最大和最小值  返回有序集合的所有元素和分数
	vals, err := pool.ZRangeByScoreWithScores("get_data", redis.ZRangeBy{
		Min: "0",
		Max: "50",
		Offset: 0,
		Count: 1,
	}).Result()

	for _, value := range vals {
		key := value.Member.(string)
		dMap, err := pool.HGetAll(key).Result()
		if err != nil {
			return url_list, err
		}
		url_list = append(url_list, dMap["url"])
		pool.Del(key).Err()
		//删除集合中的一个指定元素
		pool.ZRem("get_data", key)
	}

	return url_list, err
}

func GetKey(length int) string {
	sec := strconv.FormatInt(time.Now().Unix(), 10)
	redKey := "model_get_key:" + sec
	randLen := length
	exTime := 1
	preId := ""

	if length > 10 {
		randLen = length - 10
		preId = sec
	}
	randStr := ""
	for i := 0; i < 50; i++ {
		randStr = Random("smallnumber", randLen)
		//新增无序集合 所有的key头存在无序集合里面
		res, err := pool.SAdd(redKey, randStr, exTime).Result()
		if err == nil && res > 0 {
			break
		}
	}

	keyStr := preId + randStr
	return keyStr
}

func Random(param string, length int) string {
	str := ""
	if length < 1 {
		return str
	}
	tmp := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	switch param {
	case "number":
		tmp = "1234567890"
	case "small":
		tmp = "abcdefghijklmnopqrstuvwxyz"
	case "big":
		tmp = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "smallnumber":
		tmp = "1234567890abcdefghijklmnopqrstuvwxyz"
	case "bignumber":
		tmp = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "bigsmall":
		tmp = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	leng := len(tmp)
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		s_ind := ran.Intn(leng)
		str = str + Substr(tmp, s_ind, 1)
	}

	return str
}

/**
*  start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
*  负数 - 在从字符串结尾的指定位置开始
*  0 - 在字符串中的第一个字符处开始
*  length:正数 - 从 start 参数所在的位置返回
*  负数 - 从字符串末端返回
 */
func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	rune_str := []rune(str)
	len_str := len(rune_str)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(rune_str[start:end])
}
