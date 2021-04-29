package thread

import "sync"

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func NewWaitGroupWrapper() *WaitGroupWrapper{
    return &WaitGroupWrapper{}
}

func (w *WaitGroupWrapper) Work(wc func()){
	w.Add(1)
	go func(){
		wc()
		w.Done()
	}()

}