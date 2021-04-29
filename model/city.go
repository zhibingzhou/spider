package model

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
	err := MysqlALL["ruoai"].Create(&rc).Error
	return err
}

func GetCityFromRedis(cityName string) (map[string]string, error) {

	redisKey := "city:nick_name:" + cityName
	//优先查询redis 拿map
	dMap, err := pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var city City

		err = MysqlALL["ruoai"].Table("city").Where("title = ?", cityName).First(&city).Error
		if err == nil && city.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = city.Id
			val["province_id"] = city.Province_id
			val["title"] = city.Title

			err = pool.HMSet(redisKey, val).Err()
			if err != nil {
				return dMap, err
			}

			//新增无序集合 所有的key头存在无序集合里面
			err = pool.SAdd(head, redisKey).Err()
			if err != nil {
				return dMap, err
			}

			dMap, err = pool.HGetAll(redisKey).Result()
			if err != nil {
				return dMap, err
			}
		}

	}

	return dMap, err
}
