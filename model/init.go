package model

import "github.com/jinzhu/gorm"

var MysqlALL map[string]*gorm.DB

func init() {
	
	MysqlALL["ruoai"] = nil

	for key, _ := range MysqlALL {
		MysqlALL[key] = ReloadConfSQL("", key)
	}
}