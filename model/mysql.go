package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zhibingzhou/go_public/common"
)

var (
	//redis key
	ruoai_province = "ruoai_province"
	ruoai_city     = "ruoai_city"
)

func ReloadConfSQL(file_name string, dbname string) *gorm.DB {
	if file_name == "" {
		file_name = "./conf/database.json"
	}
	conf_byte, err := common.ReadFile(file_name)
	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	life_time, _ := time.ParseDuration(json_conf["life_time"])
	max_open, _ := strconv.Atoi(json_conf["max_open"])
	if max_open < 1 {
		max_open = 40
	}
	max_idle, _ := strconv.Atoi(json_conf["max_idle"])
	if max_idle < 1 {
		max_idle = 10
	}
	if dbname == "" {
		dbname = json_conf["db_name"]
	}
	conn_str := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8", json_conf["user"], json_conf["pwd"], json_conf["network"], json_conf["host"], json_conf["port"], dbname)
	db, err := gorm.Open("mysql", conn_str)
	if err != nil {
		fmt.Println("conn_str->", conn_str)
		panic(err)
	}
	//最大生命周期
	db.DB().SetConnMaxLifetime(life_time)
	//连接池的最大打开连接数
	db.DB().SetMaxOpenConns(max_open)
	//连接池的最大空闲连接数
	db.DB().SetMaxIdleConns(max_idle)
	db.SingularTable(true)
	//启用Logger，显示详细日志
	db.LogMode(true)

	// 禁用日志记录器，不显示任何日志
	//db.LogMode(false)
	fmt.Println("连接 mysql 成功")
	return db
}
