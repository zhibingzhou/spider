package zhenai

import (
	"fmt"

	"github.com/tebeka/selenium"
)

//定义一个任务接口，任何实现这个接口的任务，都可以用
type Job interface {
	Do(int) ParseResult
}

type ZhenAiManagerChannel struct {
	Count    int
	Request  chan Job
	Requests chan chan Job
	Result   chan ParseResult
	Finished chan int
	Driver   []selenium.WebDriver
	Service  []*selenium.Service
}

var ManagerZhenAi *ZhenAiManagerChannel

func NewZhenAiManagerChannel(count int) *ZhenAiManagerChannel {

	var zhenAi ZhenAiManagerChannel
	for i := 1; i <= count; i++ {
		driver, service := GetDriver(Port+i, SeleniumPath)
		zhenAi.Driver = append(zhenAi.Driver, driver)
		zhenAi.Service = append(zhenAi.Service, service)
	}

	return &ZhenAiManagerChannel{
		Count:    count,
		Request:  make(chan Job, count),
		Requests: make(chan chan Job, count),
		Result:   make(chan ParseResult, count),
		Finished: make(chan int, count),
		Driver:   zhenAi.Driver,
		Service:  zhenAi.Service,
	}

}

func (x *ZhenAiManagerChannel) Run() {

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

				if len(ArryReq) == 0 {
					fmt.Println("任务结束")
					x.Finished <- 1
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

//工作者拿到任务，执行
func (w Worker) Run(manager chan chan Job, re chan ParseResult, i int) {
	go func() {
		for {
			manager <- w.Work
			select {
			case work := <-w.Work:
				result := work.Do(i)
				if len(result.Requests) > 0 {
					re <- result
				}

			}
		}
	}()
}
