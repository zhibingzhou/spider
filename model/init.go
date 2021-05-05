package model

import "github.com/jinzhu/gorm"

var MysqlALL map[string]Mysql

type Mysql struct {
	DB *gorm.DB
}

func init() {
	MysqlALL = make(map[string]Mysql)
	MysqlALL["ruoai"] = Mysql{}

	for key, _ := range MysqlALL {
		db := ReloadConfSQL("", key)
		MysqlALL[key] = Mysql{DB: db}
	}
}
