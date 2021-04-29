package pro

import (
	"fmt"
)

type Pcstudy struct {
}

func (p Pcstudy) Boot() {
	fmt.Println("do something")
	Manger = NewManagerChannel(1)
	Manger.Run()
	StartChrome("https://51ao.xyz/list")

}
