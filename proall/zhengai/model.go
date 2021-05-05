package zhenai

type Request struct {
	Url      string
	PareFunc func(string, int) ParseResult
}

type ParseResult struct {
	Requests []Request
	Item     interface{}
}

func (r Request) Do(chan_number int) ParseResult {
	pareResult := r.PareFunc(r.Url, chan_number)
	return pareResult
}
