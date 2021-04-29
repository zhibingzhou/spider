package ruoai

type Request struct {
	Url        string
	ParserFunc func(string, int) ParseResult
}

type ParseResult struct {
	Requets []Request
	Items   []interface{} //列表
}

var (
	RuoaiUrl     = "https://www.ruoai.cn/location/"
	Port         = 9515
	SeleniumPath = `D:\pc\chromedriver.exe`
)


