package ruoai

import (
	"fmt"
	"test/model"
)

type Ruoai struct {
}

func (r Ruoai) Boot() {
	RuoAiManager = NewRuoAiManagerChannel(3)
	RuoAiManager.Run()
	StartRun()
}

func StartRun() {
	model.Delcash()
	RuoAiManager.Request <- Request{Url: RuoaiUrl, ParserFunc: ParseCity}

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
