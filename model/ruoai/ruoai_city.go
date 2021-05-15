package ruoai

import (
	"fmt"

	"test/model"

	"github.com/segmentio/encoding/json"
)

type City struct {
	Id          int    `json:"id" gorm:"id"`
	Province_id int    `json:"province_id" gorm:"province_id"`
	Title       string `json:"title" gorm:"title"`
}

func CityCreate(city string, province_id int) error {
	var rc City
	rc = City{
		Title:       city,
		Province_id: province_id,
	}
	err := model.MysqlALL["ruoai"].DB.Create(&rc).Error
	return err
}

func GetCityFromRedis(cityName string) (map[string]string, error) {

	redisKey := "city:nick_name:" + cityName
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var city City

		err = model.MysqlALL["ruoai"].DB.Table("city").Where("title = ?", cityName).First(&city).Error
		if err == nil && city.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = city.Id
			val["province_id"] = city.Province_id
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

	redisKey := "city:nick_name:" + fmt.Sprintf("%d", code)
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var city City

		err = model.MysqlALL["ruoai"].DB.Table("city").Where("id = ?", code).First(&city).Error
		if err == nil && city.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = city.Id
			val["province_id"] = city.Province_id
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

func GetCityProvinceFromRedis(code int) ([]City, error) {

	var city []City
	redisKey := "Province_city:nick_name:" + fmt.Sprintf("%d", code)
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["city"]) < 1 {

		err = model.MysqlALL["ruoai"].DB.Table("city").Where("province_id = ?", code).Find(&city).Error

		if err == nil && len(city) > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			array, _ := json.Marshal(city)
			val["city"] = string(array)
			err = model.Pool.HMSet(redisKey, val).Err()
			if err != nil {
				return city, err
			}

			//新增无序集合 所有的key头存在无序集合里面
			err = model.Pool.SAdd(model.Head, redisKey).Err()
			if err != nil {
				return city, err
			}

			dMap, err = model.Pool.HGetAll(redisKey).Result()
			if err != nil {
				return city, err
			}
		}

	}

	if len(dMap) > 0 {
		json.Unmarshal([]byte(dMap["city"]), &city)
	}

	return city, err
}
