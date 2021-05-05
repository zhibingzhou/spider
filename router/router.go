package router

import (
	"test/controller"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.New()
	//静态文件路径，一定需要
	Router.LoadHTMLGlob("view/*")
	Router.LoadHTMLFiles("./view/index.html", "./view/city.html", "./view/user.html", "./view/work.html", "./view/province.html", "./view/indexvue.html")
	Router.Static("/layui", "./layui")
	Router.GET("/home", controller.Home)

	//API接口 若爱
	RuoAi := Router.Group("/ruoai")
	//用户
	RuoAi.POST("/list_ruoai", controller.RuoAiList)
	//职业
	RuoAi.POST("/getEducation", controller.GetEducation)
	//城市
	RuoAi.GET("/city_list", controller.RuoAiCityList)
	//省份
	RuoAi.GET("/province_list", controller.RuoAiProvinceList)
	//省份和城市
	RuoAi.GET("/provinceCity_list", controller.RuoAiProvinceCityList)

}
