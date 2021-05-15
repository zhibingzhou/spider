package xinhe

type Film struct {
	Id          int     `json:"id"`
	En_name     string  `json:"en_name"`
	First_name  string  `json:"first_name"`
	Second_name string  `json:"second_name"`
	Year        int     `json:"year"`
	Country     string  `json:"country"`
	Score       float64 `json:"score"`
	Show_time   string  `json:"show_time"`
	Title       string  `json:"title"`
	Url_image   string  `json:"url_image"`
	Video_url   string  `json:"video_url"`
	Video_type  int  `json:"video_type"`
}

type Country struct {
	Id           int    `json:"id"`
	Country_name string `json:"country_name"`
}

type Type_all struct {
	Id        int    `json:"id"`
	Type_name string `json:"type_name"`
}

type Film_type struct {
	Id      int    `json:"id"`
	Film_id int `json:"film_id"`
	Type_id int `json:"type_id"`
}

type Film_person struct {
	Id       int    `json:"id"`
	Film_id  int    `json:"film_id"`
	Name     string `json:"name"`
	Identity int    `json:"identity"`
}
  
type Title struct {
	Id         int    `json:"id"`
	Title_name string `json:"title_name"`
}

type Film_title struct {
	Id       int    `json:"id"`
	Film_id  int `json:"film_id"`
	Title_id int `json:"title_id"`
}
