package xinhe

import "fmt"

type Job interface {
	Do() PareResult
}

type XinHeChannel struct {
	Count    int
	Request  chan Job
	Requests chan chan Job
	Result   chan PareResult
	EndJob   chan int
}

var ManagerXinHe *XinHeChannel

func NewXinHeChannel(count int) *XinHeChannel {
	return &XinHeChannel{
		Count:    count,
		Request:  make(chan Job),
		Requests: make(chan chan Job),
		Result:   make(chan PareResult),
		EndJob:   make(chan int),
	}
}

func (x *XinHeChannel) Run() {
	for i := 0; i < x.Count; i++ {
		NewWorker().Run(x.Requests, x.Result, i)
	}
	var ArryReq []Job
	var ArryReqs []chan Job

	go func() {

		for {
			var requestch chan Job
			var request Job

			if len(ArryReq) > 0 && len(ArryReqs) > 0 {
				requestch = ArryReqs[0]
				request = ArryReq[0]
				fmt.Println("通道", len(ArryReqs))
				fmt.Println("结果", len(ArryReq))
			}

			select {
			case Reqs := <-x.Requests:
				ArryReqs = append(ArryReqs, Reqs)

			case Req := <-x.Request:
				ArryReq = append(ArryReq, Req)

			case requestch <- request:
				ArryReq = ArryReq[1:]
				ArryReqs = ArryReqs[1:]

				if len(ArryReq) == 2 {
					fmt.Println("任务结束")
					x.EndJob <- 1
				}

			}

		}

	}()
}

type Worker struct {
	Work chan Job
}

func NewWorker() Worker {
	return Worker{Work: make(chan Job)}
}

func (w Worker) Run(manager chan chan Job, re chan PareResult, i int) {
	go func() {
		for {
			manager <- w.Work
			select {
			case work := <-w.Work:
				result := work.Do()
				if len(result.Requests) > 0 {
					re <- result
				}

			}
		}
	}()
}
