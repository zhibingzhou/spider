package thread

import (
	"fmt"
	"strconv"
	"strings"
	"test/model"
	"test/utils"
)

var date = "2006-01-02"
var datetime = "2006-01-02 15:04:02"

func RuoAiList(page, limit, gender, city, salary, education int, age string) (count int, rep []map[string]interface{}, err error) {

	if limit == 0 {
		limit = 10
	}
	offset := 0

	if page != 0 {
		offset = limit * (page - 1)
	}

	var users []model.User
	rep = make([]map[string]interface{}, 0)
	db := model.MysqlALL["ruoai"].DB.Table("user")

	if gender == 0 {
		db = db.Where("gender = ?", "女士")
	} else {
		db = db.Where("gender = ?", "男士")
	}

	if city != 0 {
		db = db.Where("city_code = ?", city)
	}

	if education != 0 {
		db = db.Where("education = ?", education)
	}

	if age != "" {
		age, _ := strconv.Atoi(age)
		db = db.Where("age <= ?", age)
	}
	switch salary {
	case 0:
		db = db.Where("salary_down <= ?  and  salary_up = ? ", 5000, 0)
	case 1:
		db = db.Where("salary_down >= ?  and  salary_up <= ? ", 5000, 10000)
	case 2:
		db = db.Where("salary_down >= ?  and  salary_up <= ? ", 10000, 20000)
	case 3:
		db = db.Where("salary_up >= ? ", 20000)
	default:
	}

	err = db.Limit(limit).Offset(offset).Find(&users).Error
	db.Count(&count)

	for _, value := range users {
		mapstring := utils.StructToMap(value)
		citytitle, _ := model.GetCityCodeFromRedis(value.City_code)
		mapstring["City_title"] = citytitle["title"]
		work_title, _ := model.GetWorkCodeFromRedis(value.Work)
		mapstring["Work_title"] = work_title["title"]
		hometown_title, _ := model.GetCityCodeFromRedis(value.Hometown)
		mapstring["Hometown_title"] = hometown_title["title"]
		education_title, _ := model.GetEducationCodeFromRedis(value.Education)
		mapstring["Education_title"] = education_title["title"]
		time := strings.Split(fmt.Sprintf("%s", mapstring["Birthday_time"]), " ")
		mapstring["Birthday_time"] = time[0]

		if value.Salary_down != 0 && value.Salary_up != 0 {
			mapstring["Salary"] = fmt.Sprintf("%d - %d ", value.Salary_down, value.Salary_up)
		}
		if value.Salary_down != 0 && value.Salary_up == 0 {
			mapstring["Salary"] = fmt.Sprintf("%d 以下", value.Salary_down)
		}
		if value.Salary_down == 0 && value.Salary_up != 0 {
			mapstring["Salary"] = fmt.Sprintf("%d 以上", value.Salary_up)
		}
		rep = append(rep, mapstring)
	}
	return count, rep, err

}

func RuoAiProvinceCity() (count int, rep []model.ProvinceCity, err error) {
	var pro []model.Province
	var proc []model.ProvinceCity
	err = model.MysqlALL["ruoai"].DB.Table("province").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	for _, value := range pro {
		arraycity, _ := model.GetCityProvinceFromRedis(value.Id)
		proc = append(proc, model.ProvinceCity{Title: value.Title, City: arraycity})
	}
	return len(rep), proc, err
}

func RuoAiCity() (count int, rep []map[string]interface{}, err error) {
	var pro []model.City
	err = model.MysqlALL["ruoai"].DB.Table("city").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	for _, value := range pro {
		rep = append(rep, utils.StructToMap(value))
	}
	return len(rep), rep, err
}

func RuoAiProvince() (count int, rep []model.Province, err error) {
	var pro []model.Province
	err = model.MysqlALL["ruoai"].DB.Table("province").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	return len(rep), pro, err
}

func GetEducation() (count int, rep []model.Education, err error) {
	var pro []model.Education
	err = model.MysqlALL["ruoai"].DB.Table("education").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	return len(rep), pro, err
}
