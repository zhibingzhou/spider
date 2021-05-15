package model

import (
	"github.com/jinzhu/gorm"
)

var MysqlALL map[string]Mysql

type Mysql struct {
	DB *gorm.DB
}

func init() {
	MysqlALL = make(map[string]Mysql)
	MysqlALL["ruoai"] = Mysql{}
	MysqlALL["zhenai"] = Mysql{}
	MysqlALL["xinhe"] = Mysql{}
	for key, _ := range MysqlALL {
		db := ReloadConfSQL("", key)
		MysqlALL[key] = Mysql{DB: db}
	}
}

func (m Mysql) Trans(sqlArr []string) error {
	tx := m.DB.Begin()
	// 注意，一旦你在一个事务中，使用tx作为数据库句柄
	var err error
	for _, sql := range sqlArr {
		//更新订单状态
		if err = tx.Exec(sql).Error; err != nil {
			tx.Rollback()
			break
		}
	}

	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
