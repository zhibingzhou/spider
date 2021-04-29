package xinhe

type Job interface {
	Do(int) PareResult
}

type XinHeChannel struct {
	Count    int
	Request  chan Job
	Requests chan chan Job
	Result   chan PareResult
}

var ManagerXinHe *XinHeChannel

func NewXinHeChannel(count int) *XinHeChannel {
	return &XinHeChannel{
		Count:    count,
		Request:  make(chan Job),
		Requests: make(chan chan Job),
		Result:   make(chan PareResult),
	}
}

func (x *XinHeChannel) Run() {
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

func (w Worker) Run(manager chan chan Job, re chan PareResult, i int) {
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
