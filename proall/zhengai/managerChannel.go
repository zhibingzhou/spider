package zhenai

type Job interface {
	Do(int) ParseResult
}

type ZhenAiManagerChannel struct {
	Count    int
	Request  chan Job
	Requests chan chan Job
	Result   chan ParseResult
}

var ManagerZhenAi *ZhenAiManagerChannel

func NewZhenAiManagerChannel(count int) *ZhenAiManagerChannel {
	return &ZhenAiManagerChannel{
		Count:    count,
		Request:  make(chan Job),
		Requests: make(chan chan Job),
		Result:   make(chan ParseResult),
	}
}

func (x *ZhenAiManagerChannel) Run() {
	for i := 0; i < x.Count; i++ {
		NewWorker().Run(x.Requests, x.Result, i)
	}
	go func() {
		for {
			select {
			case job := <-x.Request:
				work := <-x.Requests
				work <- job
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
