package pro

type GetData struct {
	Current_page   int    `json:"current_page"`
	Data           []Data `json:"data"`
	First_page_url string `json:"first_page_url"`
	From           int    `json:"from"`
	Last_page      int    `json:"last_page"`
	Last_page_url  string `json:"last_page_url"`
	Next_page_url  string `json:"next_page_url"`
	Path           string `json:"path"`
	Per_page       string `json:"per_page"`
	Prev_page_url  string `json:"prev_page_url"`
	To             int    `json:"to"`
	Total          int    `json:"total"`
}

type Data struct {
	Id                  int      `json:"id"`
	User_id             int      `json:"user_id"`
	Status              int      `json:"status"`
	Type                int      `json:"type"`
	Tags                []int    `json:"tags"`
	Title               string   `json:"title"`
	Is_recommend        bool     `json:"is_recommend"`
	Is_expired          bool     `json:"is_expired"`
	Source              int      `json:"source"`
	Score               int      `json:"score"`
	View_count          int      `json:"view_count"`
	City_code           string   `json:"city_code"`
	Girl_num            string   `json:"girl_num"`
	Girl_age            string   `json:"girl_age"`
	Girl_beauty         string   `json:"girl_beauty"`
	Environment         string   `json:"environment"`
	Consume_lv          string   `json:"consume_lv"`
	Serve_list          string   `json:"serve_list"`
	Serve_lv            string   `json:"serve_lv"`
	Desc                string   `json:"desc"`
	Picture             string   `json:"picture"`
	Cover_picture       string   `json:"cover_picture"`
	Anonymous           int      `json:"anonymous"`
	Published_at        string   `json:"published_at"`
	Created_at          string   `json:"created_at"`
	Updated_at          string   `json:"updated_at"`
	StatusReadable      string   `json:"statusReadable"`
	CityCodeReadable    string   `json:"cityCodeReadable"`
	UserName            string   `json:"userName"`
	UserId              int      `json:"userId"`
	UserIsStore         bool     `json:"userIsStore"`
	UserReputation      int      `json:"userReputation"`
	UserStatus          int      `json:"userStatus"`
	PublishedAtReadable string   `json:"publishedAtReadable"`
	CreatedAtReadable   string   `json:"createdAtReadable"`
	PictureHrefs        []string `json:"pictureHrefs"`
	VipProfileStatus    int      `json:"vipProfileStatus"`
	DescToHtml          string   `json:"descToHtml"`
	CoverPictureHrefs   []string `json:"coverPictureHrefs"`
}

//进程输入
type RequestCh struct {
	Data []byte
	Type int
}

//进程输出
type ReponseCh struct {
	Current_page int
	Type         int
}
