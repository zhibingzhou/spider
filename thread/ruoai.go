package thread

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"test/model"
	"test/model/ruoai"
	"test/utils"

	"github.com/olivere/elastic/v7"
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

	var users []ruoai.User
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
		citytitle, _ := ruoai.GetCityCodeFromRedis(value.City_code)
		mapstring["City_title"] = citytitle["title"]
		work_title, _ := ruoai.GetWorkCodeFromRedis(value.Work)
		mapstring["Work_title"] = work_title["title"]
		hometown_title, _ := ruoai.GetCityCodeFromRedis(value.Hometown)
		mapstring["Hometown_title"] = hometown_title["title"]
		education_title, _ := ruoai.GetEducationCodeFromRedis(value.Education)
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

func EsRuoAiList(page, limit, gender, city, salary, education int, age string) (count int, rep []map[string]interface{}, err error) {

	if limit == 0 {
		limit = 10
	}
	offset := 0

	if page != 0 {
		offset = limit * (page - 1)
	}

	var users []ruoai.User
	rep = make([]map[string]interface{}, 0)

	ctx := context.Background()
	// 创建bool查询  为组合查询
	boolQuery := elastic.NewBoolQuery().Must() //Must

	var array_query []elastic.Query

	if age != "" {
		in_age, _ := strconv.Atoi(age)
		rangeQuery := elastic.NewRangeQuery("age").
			Gt(0).     // 0 <
			Lt(in_age) // < age
		array_query = append(array_query, rangeQuery)

	}

	switch salary {
	case 0:
		termQuery := elastic.NewTermQuery("salary_up", 0)
		rangeQuery := elastic.NewRangeQuery("salary_down").
			Lte(5000) // < age
		array_query = append(array_query, rangeQuery)
		array_query = append(array_query, termQuery)

	case 1:
		rangeQuery := elastic.NewRangeQuery("salary_up").
			Lte(10000) // < age
		rangeQuery1 := elastic.NewRangeQuery("salary_down").
			Gte(5000) // < age
		array_query = append(array_query, rangeQuery)
		array_query = append(array_query, rangeQuery1)
	case 2:
		rangeQuery := elastic.NewRangeQuery("salary_up").
			Lte(20000) // < age
		rangeQuery1 := elastic.NewRangeQuery("salary_down").
			Gte(10000) // < age
		array_query = append(array_query, rangeQuery)
		array_query = append(array_query, rangeQuery1)
	case 3:
		rangeQuery := elastic.NewRangeQuery("salary_up").
			Gte(20000) // < age
		array_query = append(array_query, rangeQuery)
	default:
	}

	if gender == 0 {

		// 创建term查询条件，用于精确查询
		termQuery := elastic.NewTermQuery("gender", "女")
		termQuery1 := elastic.NewTermQuery("gender", "士")
		array_query = append(array_query, termQuery)
		array_query = append(array_query, termQuery1)

	} else {
		termQuery := elastic.NewTermQuery("gender", "男")
		termQuery1 := elastic.NewTermQuery("gender", "士")
		array_query = append(array_query, termQuery)
		array_query = append(array_query, termQuery1)
	}

	if city != 0 {
		termQuery := elastic.NewTermQuery("city_code", city)
		array_query = append(array_query, termQuery)
	}

	if education != 0 {
		termQuery := elastic.NewTermQuery("education", education)
		array_query = append(array_query, termQuery)
	}

	boolQuery.Must(array_query...)

	searchResult, err := model.EsRuoai.ElaClient.Search().
		Index(model.EsRuoai.Index). // 设置索引名
		Query(boolQuery).
		Sort("id", true).     // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(offset).         // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(limit).          // 设置分页参数 - 每页大小
		Pretty(true).         // 查询结果返回可读性较好的JSON格式
		TrackTotalHits(true). //可以查所有条数
		Do(ctx)               // 执行请求

	if err != nil {
		// Handle error
		panic(err)
	}

	if searchResult.TotalHits() > 0 {
		// 查询结果不为空，则遍历结果
		var b1 ruoai.User
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
			// 转换成User对象
			if t, ok := item.(ruoai.User); ok {
				fmt.Println(t)
				users = append(users, t)
			}
		}

	}
	//  searchResult.TookInMillis  == 上面拿 count 的方法
	fmt.Printf("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())

	count = int(searchResult.TotalHits())

	for _, value := range users {
		mapstring := utils.StructToMap(value)
		citytitle, _ := ruoai.GetCityCodeFromRedis(value.City_code)
		mapstring["City_title"] = citytitle["title"]
		work_title, _ := ruoai.GetWorkCodeFromRedis(value.Work)
		mapstring["Work_title"] = work_title["title"]
		hometown_title, _ := ruoai.GetCityCodeFromRedis(value.Hometown)
		mapstring["Hometown_title"] = hometown_title["title"]
		education_title, _ := ruoai.GetEducationCodeFromRedis(value.Education)
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

func ShowSelect(searchResult *elastic.SearchResult) {
	if searchResult.TotalHits() > 0 {
		// 查询结果不为空，则遍历结果
		var b1 ruoai.User
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
			// 转换成Article对象
			if t, ok := item.(ruoai.User); ok {
				fmt.Println(t)
			}
		}

	}
}

func RuoAiProvinceCity() (count int, rep []ruoai.ProvinceCity, err error) {
	var pro []ruoai.Province
	var proc []ruoai.ProvinceCity
	err = model.MysqlALL["ruoai"].DB.Table("province").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	for _, value := range pro {
		arraycity, _ := ruoai.GetCityProvinceFromRedis(value.Id)
		proc = append(proc, ruoai.ProvinceCity{Title: value.Title, City: arraycity})
	}
	return len(rep), proc, err
}

func RuoAiCity() (count int, rep []map[string]interface{}, err error) {
	var pro []ruoai.City
	err = model.MysqlALL["ruoai"].DB.Table("city").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	for _, value := range pro {
		rep = append(rep, utils.StructToMap(value))
	}
	return len(rep), rep, err
}

func RuoAiProvince() (count int, rep []ruoai.Province, err error) {
	var pro []ruoai.Province
	err = model.MysqlALL["ruoai"].DB.Table("province").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	return len(rep), pro, err
}

func GetEducation() (count int, rep []ruoai.Education, err error) {
	var pro []ruoai.Education
	err = model.MysqlALL["ruoai"].DB.Table("education").Find(&pro).Error
	if err != nil {
		return 0, rep, err
	}
	return len(rep), pro, err
}
