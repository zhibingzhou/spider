package zhenai

import (
	"fmt"
	"test/model/zhenai"
	"test/utils"
)

type ZhenAi struct {
	Url string
}

var ZhenAiT ZhenAi

func (z ZhenAi) Boot() {
    
    ManagerZhenAi = NewZhenAiManagerChannel(1)

	PareUserInformation("https://album.zhenai.com/u/1095837354",0)
    //ManagerZhenAi.Run()
	//StartRun(z.Url)
}

func init() {
	ZhenAiT = ZhenAi{Url: "https://www.zhenai.com/zhenghun"}
}

func StartRun(Url string) {
	//model.Delcash()

	GetCity(Url)

	go func() {
		for {

			value, err := zhenai.GetCityData()
			if err != nil {
				fmt.Println("StartRun", err)
				continue
			}

			if len(value) < 1 {
				fmt.Println("未拿到城市--说明爬完", err)
				return
			}

			for _, val := range value {
				if val != "" {
					fmt.Println("拿到城市url ->", val)
					ManagerZhenAi.Request <- Request{Url: val, PareFunc: PareAllPage}
					//RuoAiManager.Request <- Request{Url: val, ParserFunc: GetBoyUrl}
				}
			}
			//等10秒
			utils.RandSleep(10)
			<-ManagerZhenAi.Finished
			fmt.Println("拿到任务")
		}

	}()

	go func() {
		for {
			result := <-ManagerZhenAi.Result
			fmt.Println("拿到请求有", len(result.Requests))
			if len(result.Requests) < 1 {
				continue
			}
			for _, value := range result.Requests {
				ManagerZhenAi.Request <- value
			}
		}
	}()

	select {}

}
