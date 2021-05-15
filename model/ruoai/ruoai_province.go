package ruoai

import "test/model"

type Province struct {
	Id    int    `json:"id" gorm:"id"`
	Title string `json:"title" gorm:"title"`
}

type ProvinceCity struct {
	Title string
	City  []City
}

func ProvinceCreate(province string) error {
	var rc Province
	rc = Province{
		Title: province,
	}
	err := model.MysqlALL["ruoai"].DB.Create(&rc).Error
	return err
}

func GetProvinceFromRedis(proName string) (map[string]string, error) {

	redisKey := "province:nick_name:" + proName
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var province Province

		err = model.MysqlALL["ruoai"].DB.Table("province").Where("title = ?", proName).First(&province).Error
		if err == nil && province.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = province.Id
			val["title"] = province.Title

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
