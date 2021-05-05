package ruoai

import (
	"fmt"
	"test/model"
	"test/utils"
)

type Ruoai struct {
}

func (r Ruoai) Boot() {
	RuoAiManager = NewRuoAiManagerChannel(10)
	RuoAiManager.Run()
	StartRun()
}

func StartRun() {
	//model.Delcash()
	//RuoAiManager.Request <- Request{Url: RuoaiUrl, ParserFunc: GetBoyUrl}

	ParseCity(RuoaiUrl, 1)

	go func() {
		for {

			value, err := model.GetCityData()
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
					RuoAiManager.Request <- Request{Url: val, ParserFunc: GetGirlUrl}
					RuoAiManager.Request <- Request{Url: val, ParserFunc: GetBoyUrl}
				}
			}
			//等10秒
			utils.RandSleep(10)
			<-RuoAiManager.WriteMysql
			fmt.Println("拿到任务")
		}

	}()

	go func() {
		for {
			result := <-RuoAiManager.Result
			fmt.Println("拿到请求有", len(result.Requets))
			if len(result.Requets) < 1 {
				continue
			}
			for _, value := range result.Requets {
				RuoAiManager.Request <- value
			}
		}
	}()

	select {}

}
