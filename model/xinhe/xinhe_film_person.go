package xinhe

import (
	"fmt"
	"test/model"
)

func FilmPersonCreate(film_id, Identity int, name string) error {
	var rc Film_person
	rc = Film_person{
		Film_id:  film_id,
		Name:     name,
		Identity: Identity,
	}
	err := model.MysqlALL["xinhe"].DB.Create(&rc).Error
	return err
}

func GetFilmPersonFromRedis(id, identity int, name string) (map[string]string, error) {

	redisKey := "film_person:nick_name:" + fmt.Sprintf("%d_%d_%s", id, identity, name)
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var film_person Film_person

		err = model.MysqlALL["xinhe"].DB.Table("film_person").Where("identity = ? and film_id = ? and name = ?", identity, id, name).First(&film_person).Error
		if err == nil && film_person.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = film_person.Id
			val["identity"] = film_person.Identity
			val["name"] = film_person.Name

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
