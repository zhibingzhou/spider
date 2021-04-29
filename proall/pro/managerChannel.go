package pro

import (
	"encoding/json"
	"fmt"
)

var Manger *ManagerChannel

type ManagerChannel struct {
	Count    int
	Request  chan RequestCh
	Requests chan chan RequestCh
	Reponse  chan ReponseCh
}

func NewManagerChannel(count int) *ManagerChannel {
	return &ManagerChannel{
		Count:    count,
		Request:  make(chan RequestCh),
		Requests: make(chan chan RequestCh),
		Reponse:  make(chan ReponseCh),
	}
}

type Worker struct {
	Wrequest chan RequestCh
}

func NewWoker() Worker {
	return Worker{Wrequest: make(chan RequestCh)}
}

func (m ManagerChannel) Run() {
	//先让每个工作都等带任务
	for i := 0; i < m.Count; i++ {
		NewWoker().ChannelRun(m.Requests, m.Reponse)
	}
	var ArryReq []RequestCh
	var ArryReqs []chan RequestCh

	go func() {

		for {
			var requestch chan RequestCh
			var request RequestCh

			if len(ArryReq) > 0 && len(ArryReqs) > 0 {
				requestch = ArryReqs[0]
				request = ArryReq[0]
			}

			select {
			case Reqs := <-m.Requests:
				ArryReqs = append(ArryReqs, Reqs)

			case Req := <-m.Request:
				ArryReq = append(ArryReq, Req)

			case requestch <- request:
				ArryReq = ArryReq[1:]
				ArryReqs = ArryReqs[1:]
			}
		}

	}()

}

//每个工作者
func (w Worker) ChannelRun(crequest chan chan RequestCh, req chan ReponseCh) {
	//

	go func() {
		for {
			crequest <- w.Wrequest //这个 是告诉管理员，我准备好了
			select {
			case re := <-w.Wrequest: //等待管理员分配任务
				result := re.Run()
				req <- result
			}
		}
	}()
}

func (r RequestCh) Run() ReponseCh {
	var getData GetData
	// re, err := utils.ZhToUnicode(string(r.Data))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//res := strings.Replace(string(re), `\`, "", -1)
	err := json.Unmarshal([]byte(r.Data), &getData)
	if err != nil {
		fmt.Println(err)
	}
	return ReponseCh{Current_page: getData.Current_page, Type: r.Type}
}
