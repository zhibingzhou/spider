package thread

import (
	"fmt"
	"strings"
	"test/model"
	"test/model/xinhe"
	"test/utils"
)

func GetFilm(film_type, page, limit int) (count int, rep []map[string]interface{}, err error) {
	var film []xinhe.Film
	rep = make([]map[string]interface{}, 0)
	if limit == 0 {
		limit = 10
	}

	offset := 0

	if page != 0 {
		offset = limit * (page - 1)
	}

	seletct_sql := "SELECT film.`id`,film.`first_name`,film.`score`,film.`second_name`,film.`show_time`,film.`title`,film.`year`,film.`video_type`,film.`video_url`,film.`url_image`,film.`en_name`, LEFT(film.`country`,LENGTH(film.`country`)-1) AS country"
	select_count := "SELECT COUNT(*) AS COUNT "
	tableName := " from film "
	where_sql := "where video_url != '' "

	// if year != "" {
	// 	where_sql = where_sql + " and year = " + year
	// }

	// if title_id != "" {
	// 	where_sql = where_sql + " and title_id like '%" + title_id + "%'"
	// }

	// if person_name != "" {
	// 	where_sql = where_sql + " and person_name = " + person_name
	// }

	if film_type != 0 {
		tableName = tableName + " left join film_type on  film_type.film_id = film.id  "
		where_sql = where_sql + " and film_type.type_id = " + fmt.Sprintf("%d", film_type)
	}

	err = model.MysqlALL["xinhe"].DB.Debug().Raw(select_count + tableName + where_sql).Count(&count).Error
	
	where_sql = where_sql + fmt.Sprintf(" order by YEAR desc limit %d,%d;", offset, limit)

	err = model.MysqlALL["xinhe"].DB.Debug().Raw(seletct_sql + tableName + where_sql).Scan(&film).Error

	for _, value := range film {
		country_array := strings.Split(value.Country, ",")

		pro := utils.StructToMap(value)
		var country_list []string
		for _, value1 := range country_array {
			country_name, _ := xinhe.GetCountryByIdRedis(value1)
			country_list = append(country_list, country_name["country_name"])
		}
		delete(pro, "country_name")
		pro["country_name"] = country_list
		rep = append(rep, pro)
	}

	if err != nil {
		return 0, rep, err
	}
	return count, rep, err
}

func XinHeType() (count int, rep []xinhe.Type_all, err error) {
	db := model.MysqlALL["xinhe"].DB.Table("type_all")
	err = db.Find(&rep).Error
	db.Count(&count)
	return count, rep, err
}

func GetFilmByYear(year, title_id, person_name string) string {
	where_sql := "where year = " + year

	if title_id != "" {
		where_sql = where_sql + " LEFT JOIN film_title ON film.`id` = film_title.`film_id` WHERE film_title.`title_id` IN (2186,1233) "
	}

	if person_name != "" {
		where_sql = where_sql + " and person_name = " + person_name
	}

	return ""
}
