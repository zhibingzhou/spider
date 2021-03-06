package ruoai

import "test/model"

type User struct {
	Id            int    `json:"id" gorm:"id"`                       //
	Nickname      string `json:"nickname" gorm:"nickname"`           //
	City_code     int    `json:"city_code" gorm:"city_code"`         //
	Age           int    `json:"age" gorm:"age"`                     //
	Gender        string `json:"gender" gorm:"gender"`               //
	Education     int    `json:"education" gorm:"education"`         //
	Birthday_time string `json:"birthday_time" gorm:"birthday_time"` //
	Height        string `json:"height" gorm:"height"`               //
	Weight        string `json:"weight" gorm:"weight"`               //
	Married       string `json:"married" gorm:"married"`             //
	House         string `json:"house" gorm:"house"`                 //
	Work          int    `json:"work" gorm:"work"`                   //
	Hometown      int    `json:"hometown" gorm:"hometown"`           //
	Salary_up     int    `json:"salary_up" gorm:"salary_up"`         //
	Salary_down   int    `json:"salary_down" gorm:"salary_down"`     //
	Title         string `json:"title" gorm:"title"`                 //
	Url_image     string `json:"url_image" gorm:"url_image"`         //
}

type ReponseUser struct {
	Id              int    `json:"id" gorm:"id"`                           //
	Nickname        string `json:"nickname" gorm:"nickname"`               //
	City_title      int    `json:"city_title" gorm:"city_title"`           //
	Age             int    `json:"age" gorm:"age"`                         //
	Gender          string `json:"gender" gorm:"gender"`                   //
	Education_title string `json:"education_title" gorm:"education_title"` //
	Birthday_time   string `json:"birthday_time" gorm:"birthday_time"`     //
	Height          string `json:"height" gorm:"height"`                   //
	Weight          string `json:"weight" gorm:"weight"`                   //
	Married         string `json:"married" gorm:"married"`                 //
	House           string `json:"house" gorm:"house"`                     //
	Work_title      string `json:"work_title" gorm:"work_title"`           //
	Hometown_title  string `json:"hometown_title" gorm:"hometown_title"`   //
	Salary_up       int    `json:"salary_up" gorm:"salary_up"`             //
	Salary_down     int    `json:"salary_down" gorm:"salary_down"`         //
	Title           string `json:"title" gorm:"title"`                     //
	Url_image       string `json:"url_image" gorm:"url_image"`             //
}

func UserCreate(u User) error {
	err := model.MysqlALL["ruoai"].DB.Create(&u).Error
	return err
}

func GetUserFromRedis(id string) (map[string]string, error) {

	redisKey := "user:nick_name:" + id
	//????????????redis ???map
	dMap, err := model.Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var user User

		err = model.MysqlALL["ruoai"].DB.Table("user").Where("id = ?", id).First(&user).Error
		if err == nil && user.Id > 0 {
			// ??????????????? ??? map
			val := map[string]interface{}{}
			val["id"] = user.Id
			err = model.Pool.HMSet(redisKey, val).Err()
			if err != nil {
				return dMap, err
			}

			//?????????????????? ?????????key???????????????????????????
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
