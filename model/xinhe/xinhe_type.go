package xinhe

import "test/model"

func TypeCreate(film_type string) error {
	var rc Type_all
	rc = Type_all{
		Type_name: film_type,
	}
	err := model.MysqlALL["xinhe"].DB.Create(&rc).Error
	return err
}

func GetFilmTypeFromRedis(proName string) (map[string]string, error) {

	redisKey := "film_type:nick_name:" + proName
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var film_type Type_all

		err = model.MysqlALL["xinhe"].DB.Table("type_all").Where("type_name = ?", proName).First(&film_type).Error
		if err == nil && film_type.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = film_type.Id
			val["type_name"] = film_type.Type_name

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

func GetTypeNameByIdRedis(Id string) (map[string]string, error) {

	redisKey := "film_type:nick_name:" + Id
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var film_type Type_all

		err = model.MysqlALL["xinhe"].DB.Table("type_all").Where("id = ?", Id).First(&film_type).Error
		if err == nil && film_type.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = film_type.Id
			val["type_name"] = film_type.Type_name

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
