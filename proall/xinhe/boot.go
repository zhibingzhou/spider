package xinhe

import (
	"fmt"
	"test/model/xinhe"
)

type XinHe struct {
}

func (x XinHe) Boot() {

	//写入任务,爬电影
	GetMovieStart("https://xinghe.tv/movie")

	ManagerXinHe = NewXinHeChannel(2)
	ManagerXinHe.Run()

	go func() {
		for {
			result := <-ManagerXinHe.Result
			fmt.Println("拿到请求有", len(result.Requests))
			if len(result.Requests) < 1 {
				continue
			}
			for _, value := range result.Requests {
				ManagerXinHe.Request <- value
			}
		}
	}()

	for {
		xh, err := xinhe.GettypeData()
		if err != nil || len(xh) < 1 {
			fmt.Println("结束任务")
			break
		}

		for _, value := range xh {
			if value != "" {
				ManagerXinHe.Request <- Request{Url: value, PareFunc: FindScr, Type: 1}
			}
		}

		<-ManagerXinHe.EndJob
		fmt.Println("再取类型")
	}

	select {}
}
