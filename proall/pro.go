package proall

import (
	"fmt"
	"os"
	"test/proall/pro"
	"test/proall/ruoai"
	"test/proall/xinhe"
	"test/thread"
	"test/utils"
)

var (
	Ctn int = 0
)

const (
	SCAN_TASK_MAX_NUM      = 0
	GCRAWLER_SCAN_MAX_RATE = 20
	GCRAWLER_GRAB_MAX_RATE = 2
	HUOLI_TABLE_ID         = 1
	JRSKQ_TABLE_ID         = 5
)

var (
	pros   pro.Pcstudy
	ruoais ruoai.Ruoai
	xinhes xinhe.XinHe
)

//结构
type Processor struct {
	Calls map[string]func()
}

//申请空间
func NewProcessor() *Processor {
	return &Processor{Calls: make(map[string]func())}
}

//注册
func (this *Processor) Register() *Processor {

	this.Calls = map[string]func(){
		//"51": pros.Boot,
		//"ruoai": ruoais.Boot,
		//"zhenai": zhenai.ZhenAiT.Boot,
		"xinhe": xinhes.Boot,
	}
	return this
}

func (this *Processor) Boot() {
	for {
		Wall := thread.NewWaitGroupWrapper()
		for _, value := range this.Calls {
			Wall.Work(value)
		}
		utils.RandSleep(GCRAWLER_SCAN_MAX_RATE)
		Wall.Wait()
		this.Exit()
	}
}

//Exit
func (this *Processor) Exit() {
	if Ctn == SCAN_TASK_MAX_NUM {
		fmt.Println("执行次数！")
		os.Exit(0)
	}
	Ctn++
}
