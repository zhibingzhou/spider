package controller

import "github.com/gin-gonic/gin"

/*
* 定义一个json的返回值类型(首字母必须大写)
 */
type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

func Ok(data interface{}, count int, c *gin.Context) {
	c.JSON(http_status, &Response{Code: 200, Msg: "success", Data: data, Count: count})
}

func Fail(data interface{}, c *gin.Context) {
	c.JSON(http_status, &Response{Code: 100, Msg: "fail", Data: data})
}
