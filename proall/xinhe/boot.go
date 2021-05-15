package xinhe

import (
	"fmt"
)

type XinHe struct {
}

func (x XinHe) Boot() {

	// GetVideo("https://xinghe.tv/play/26979204", 1)
	//写入任务,爬电影

	//GetMovieStart("https://xinghe.tv/movie")

	ManagerXinHe = NewXinHeChannel(4)
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
		// xh, err := xinhe.GettypeData()
		// if err != nil || len(xh) < 1 {
		// 	fmt.Println("结束任务")
		// 	break
		// }
		xh := []string{"/movie?type=movie&genre=%E5%85%A8%E9%83%A8&keyword=&region=all&sort=hot&year=2021"}
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
