package ruoai

import (
	"fmt"
	"test/model"
)

type Education struct {
	Id    int    `json:"id" gorm:"id"`
	Title string `json:"title" gorm:"title"`
}

func EducationCreate(education string) error {
	var rc Education
	rc = Education{
		Title: education,
	}
	err := model.MysqlALL["ruoai"].DB.Create(&rc).Error
	return err
}

func GetEducationFromRedis(proName string) (map[string]string, error) {

	redisKey := "education:nick_name:" + proName
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var education Education

		err = model.MysqlALL["ruoai"].DB.Table("education").Where("title = ?", proName).First(&education).Error
		if err == nil && education.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = education.Id
			val["title"] = education.Title

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

func GetEducationCodeFromRedis(proName int) (map[string]string, error) {

	redisKey := "education:nick_name:" + fmt.Sprintf("%d", proName)
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var education Education

		err = model.MysqlALL["ruoai"].DB.Table("education").Where("id = ?", proName).First(&education).Error
		if err == nil && education.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = education.Id
			val["title"] = education.Title

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
