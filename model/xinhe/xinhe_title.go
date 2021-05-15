package xinhe

import "test/model"

func TitleCreate(title_name string) error {
	var rc Title
	rc = Title{
		Title_name: title_name,
	}
	err := model.MysqlALL["xinhe"].DB.Create(&rc).Error
	return err
}

func GetFilmTitleFromRedis(proName string) (map[string]string, error) {

	redisKey := "title_name:nick_name:" + proName
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var title_name Title

		err = model.MysqlALL["xinhe"].DB.Table("title").Where("title_name = ?", proName).First(&title_name).Error
		if err == nil && title_name.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = title_name.Id
			val["title_name"] = title_name.Title_name

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

func GetTitleByIdRedis(Id string) (map[string]string, error) {

	redisKey := "title_name:nick_name:" + Id
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var title_name Title

		err = model.MysqlALL["xinhe"].DB.Table("title").Where("id = ?", Id).First(&title_name).Error
		if err == nil && title_name.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = title_name.Id
			val["title_name"] = title_name.Title_name

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
