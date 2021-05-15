package xinhe

type Request struct {
	Url      string
	Type     int
	PareFunc func(string, int) PareResult
}

type PareResult struct {
	Requests []Request
	Item     interface{}
}

func (r Request) Do() PareResult {
	pareResult := r.PareFunc("https://xinghe.tv"+r.Url, r.Type)
	return pareResult
}
