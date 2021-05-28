package controller

import (
	"test/thread"

	"github.com/gin-gonic/gin"
)

type XinHeRequestList struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	Film_type int `json:"film_type"`
}

func XinHeList(c *gin.Context) {
	var xinhe XinHeRequestList
	c.ShouldBindJSON(&xinhe)

	count, d, err := thread.GetFilm(xinhe.Film_type, xinhe.Page, xinhe.PageSize)
	if err != nil {
		Fail(d, c)
		return
	}
	Ok(d, count, c)
}

func XinHeType(c *gin.Context) {
	count, d, err := thread.XinHeType()
	if err != nil {
		Fail(d, c)
		return
	}
	Ok(d, count, c)
}

func Video(c *gin.Context) {
	c.HTML(200, "videoindex.html", gin.H{
		"title": "新河影院",
	})
}
