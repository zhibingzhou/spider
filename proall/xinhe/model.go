package xinhe

type Request struct {
	Url      string
	PareFunc func(string, int) PareResult
}

type PareResult struct {
	Requests []Request
	Item     interface{}
}

func (r Request) Do(chan_number int) PareResult {
	pareResult := r.PareFunc(r.Url, chan_number)
	return pareResult
}
