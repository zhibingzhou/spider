package ruoai

import (
	"fmt"
	"test/model"
)

type RuoAiManagerChannel struct {
	Count      int
	Request    chan Request
	Requests   chan chan Request
	Result     chan ParseResult
	WriteMysql chan model.User
}

type Wokers struct {
	Work chan Request
}

var RuoAiManager *RuoAiManagerChannel

func NewRuoAiManagerChannel(count int) *RuoAiManagerChannel {
	return &RuoAiManagerChannel{
		Count:      count,
		Request:    make(chan Request),
		Requests:   make(chan chan Request),
		Result:     make(chan ParseResult),
		WriteMysql: make(chan model.User),
	}
}

func (r *RuoAiManagerChannel) Run() {
	for i := 0; i < r.Count; i++ {
		NewWoker().Run(r.Requests, r.Result, i)
	}

	var ArryReq []Request
	var ArryReqs []chan Request

	go func() {

		for {
			var requestch chan Request
			var request Request

			if len(ArryReq) > 0 && len(ArryReqs) > 0 {
				requestch = ArryReqs[0]
				request = ArryReq[0]
				fmt.Println("通道", len(ArryReqs))
				fmt.Println("结果", len(ArryReq))
			}

			select {
			case Reqs := <-r.Requests:
				ArryReqs = append(ArryReqs, Reqs)

			case Req := <-r.Request:
				ArryReq = append(ArryReq, Req)

			case requestch <- request:
				ArryReq = ArryReq[1:]
				ArryReqs = ArryReqs[1:]
			}
		}

	}()
}

func NewWoker() Wokers {
	return Wokers{
		Work: make(chan Request),
	}
}

func (w Wokers) Run(r chan chan Request, p chan ParseResult, count int) {
	go func() {
		for {
			r <- w.Work
			select {
			case wd := <-w.Work:
				fmt.Println("执行")
				result := wd.ParserFunc(wd.Url, count)
				if len(result.Requets) > 0 {
					p <- result
				}
			}
		}
	}()
}
