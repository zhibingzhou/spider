package xinhe

import "test/model"

func CountryCreate(country string) error {
	var rc Country
	rc = Country{
		Country_name: country,
	}
	err := model.MysqlALL["xinhe"].DB.Create(&rc).Error
	return err
}

func GetCountryFromRedis(proName string) (map[string]string, error) {

	redisKey := "country:nick_name:" + proName
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var country Country

		err = model.MysqlALL["xinhe"].DB.Table("country").Where("country_name = ?", proName).First(&country).Error
		if err == nil && country.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = country.Id
			val["country_name"] = country.Country_name

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

func GetCountryByIdRedis(Id string) (map[string]string, error) {

	redisKey := "country:nick_name:" + Id
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var country Country

		err = model.MysqlALL["xinhe"].DB.Table("country").Where("id = ?", Id).First(&country).Error
		if err == nil && country.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = country.Id
			val["country_name"] = country.Country_name

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
