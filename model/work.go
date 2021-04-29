package model

type Work struct {
	Id    int    `json:"id" gorm:"id"`
	Title string `json:"title" gorm:"title"`
}

func WorkCreate(work string) error {
	var rc Work
	rc = Work{
		Title: work,
	}
	err := MysqlALL["ruoai"].Create(&rc).Error
	return err
}

func GetWorkFromRedis(proName string) (map[string]string, error) {

	redisKey := "work:nick_name:" + proName
	//优先查询redis 拿map
	dMap, err := pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var work Work

		err = MysqlALL["ruoai"].Table("work").Where("title = ?", proName).First(&work).Error
		if err == nil && work.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = work.Id
			val["title"] = work.Title

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
