package xinhe

import "test/model"

func FilmTypeCreate(film_id, type_id int) error {
	var rc Film_type
	rc = Film_type{
		Film_id: film_id,
		Type_id: type_id,
	}
	err := model.MysqlALL["xinhe"].DB.Create(&rc).Error
	return err
}
