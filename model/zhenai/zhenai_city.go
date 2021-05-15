package zhenai

import (
	"fmt"

	"test/model"
)

type City struct {
	Id    int    `json:"id" gorm:"id"`
	Code  string `json:"code" gorm:"code"`
	Title string `json:"title" gorm:"title"`
}

func CityCreate(city, code string) error {
	var rc City
	rc = City{
		Title: city,
		Code:  code,
	}
	err := model.MysqlALL["zhenai"].DB.Create(&rc).Error
	return err
}

func GetCityFromRedis(cityName string) (map[string]string, error) {

	redisKey := "zhenai_city:nick_name:" + cityName
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var city City

		err = model.MysqlALL["zhenai"].DB.Table("City").Where("title = ?", cityName).First(&city).Error
		if err == nil && city.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = city.Id
			val["code"] = city.Code
			val["title"] = city.Title

			err = model.Pool.HMSet(redisKey, val).Err()
			if err != nil {
				return dMap, err
			}

			//新增无序集合 所有的key头存在无序集合里面
			err = model.Pool.SAdd(model.Head, redisKey).Err()
			if err != nil {
				return dMap, err
			}

			dMap, err = model.Pool.HGetAll(redisKey).Result()
			if err != nil {
				return dMap, err
			}
		}

	}

	return dMap, err
}

func GetCityCodeFromRedis(code int) (map[string]string, error) {

	redisKey := "zhenai_city:nick_name:" + fmt.Sprintf("%d", code)
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var city City

		err = model.MysqlALL["zhenai"].DB.Table("city").Where("id = ?", code).First(&city).Error
		if err == nil && city.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = city.Id
			val["code"] = city.Code
			val["title"] = city.Title

			err = model.Pool.HMSet(redisKey, val).Err()
			if err != nil {
				return dMap, err
			}

			//新增无序集合 所有的key头存在无序集合里面
			err = model.Pool.SAdd(model.Head, redisKey).Err()
			if err != nil {
				return dMap, err
			}

			dMap, err = model.Pool.HGetAll(redisKey).Result()
			if err != nil {
				return dMap, err
			}
		}

	}

	return dMap, err
}
