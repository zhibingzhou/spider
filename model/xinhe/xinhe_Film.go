package xinhe

import (
	"test/model"
)

type FilmReponse struct {
	Id           int           `json:"id"`
	En_name      string        `json:"en_name"`
	First_name   string        `json:"first_name"`
	Second_name  string        `json:"second_name"`
	Year         int           `json:"year"`
	Country      string        `json:"country"`
	Score        float64       `json:"score"`
	Show_time    string        `json:"show_time"`
	Title        string        `json:"title"`
	Url_image    string        `json:"url_image"`
	Video_url    string        `json:"video_url"`
	Video_type   int           `json:"video_type"`
	Country_name []string      `json:"country_name"`
	Title_name   []Title       `json:"title_name"`
	Person       []Film_person `json:"person"`
}

func FilmCreate(rc Film) error {
	err := model.MysqlALL["xinhe"].DB.Create(&rc).Error
	return err
}

func UpdateFilm(id int, m map[string]interface{}) error {
	err := model.MysqlALL["xinhe"].DB.Table("film").Where(" id = ? ", id).UpdateColumns(m).Error
	return err
}

func GetFilmByIdRedis(Id string) (map[string]string, error) {

	redisKey := "film:nick_name:" + Id
	//优先查询redis 拿map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var film Title

		err = model.MysqlALL["xinhe"].DB.Table("film").Where("id = ?", Id).First(&film).Error
		if err == nil && film.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = film.Id

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

func FilmList() {

}
