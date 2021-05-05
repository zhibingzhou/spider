package controller

import (
	"test/thread"

	"github.com/gin-gonic/gin"
)

var http_status = 200

type RuoaiReponseList struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	Sex       int `json:"sex"`
	City      int `json:"city"`
	Salary    int `json:"salary"`
	Education int `json:"education"`
	Age       string `json:"age"`
}

func Home(c *gin.Context) {
	_, d, _ := thread.RuoAiProvinceCity()
	c.HTML(200, "indexvue.html", gin.H{
		"title": "若爱",
		"all":   d,
	})
}

func HomeVue(c *gin.Context) {
	c.HTML(200, "vue.html", gin.H{
		"title": "若爱",
	})
}

func RuoAiList(c *gin.Context) {
	var ruoai RuoaiReponseList
	_ = c.ShouldBindJSON(&ruoai)

	count, d, err := thread.RuoAiList(ruoai.Page, ruoai.PageSize, ruoai.Sex, ruoai.City, ruoai.Salary, ruoai.Education ,ruoai.Age )
	if err != nil {
		Fail(d, c)
		return
	}
	Ok(d, count, c)
}

func GetEducation(c *gin.Context) {

	count, d, err := thread.GetEducation()
	if err != nil {
		Fail(d, c)
		return
	}
	Ok(d, count, c)
}

func RuoAiCityList(c *gin.Context) {
	count, d, err := thread.RuoAiCity()
	if err != nil {
		Fail(d, c)
		return
	}
	Ok(d, count, c)
}

func RuoAiProvinceList(c *gin.Context) {
	count, d, err := thread.RuoAiProvince()
	if err != nil {
		Fail(d, c)
		return
	}
	Ok(d, count, c)
}

func RuoAiProvinceCityList(c *gin.Context) {
	count, d, err := thread.RuoAiProvinceCity()
	if err != nil {
		Fail(d, c)
		return
	}
	Ok(d, count, c)
}
