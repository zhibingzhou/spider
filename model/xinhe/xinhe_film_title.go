package xinhe

import (
	"fmt"
	"test/model"
)

func FilmTitleCreate(film_id, title_id int) error {
	var rc Film_title
	rc = Film_title{
		Film_id:  film_id,
		Title_id: title_id,
	}
	err := model.MysqlALL["xinhe"].DB.Create(&rc).Error
	return err
}

func PGetFilmTitleFromRedis(film_id, title_id int) (map[string]string, error) {

	redisKey := "film_title:nick_name:" + fmt.Sprintf("%d_%d", film_id, title_id)
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var film_title Film_title

		err = model.MysqlALL["xinhe"].DB.Table("film_title").Where("film_id = ? and title_id = ?", film_id, title_id).First(&film_title).Error
		if err == nil && film_title.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = film_title.Id
			val["film_id"] = film_title.Film_id
			val["title_id"] = film_title.Title_id

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
