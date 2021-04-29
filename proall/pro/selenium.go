package pro

import (
	"fmt"
	"test/utils"
	"time"

	"github.com/tebeka/selenium"
)

var GetPostData = `
(function(XHR) {
  "use strict";

  var element = document.createElement('textarea');

  element.id = "interceptedResponse";
  element.appendChild(document.createTextNode(""));

  document.body.appendChild(element);



  var open = XHR.prototype.open;
  var send = XHR.prototype.send;

  XHR.prototype.open = function(method, url, async, user, pass) {
    this._url = url; // want to track the url requested
    open.call(this, method, url, async, user, pass);
  };

  XHR.prototype.send = function(data) {
    var self = this;
    var oldOnReadyStateChange;
    var url = this._url;

    function onReadyStateChange() {
      if(self.status === 200 && self.readyState == 4 /* complete */) {
        if( self.responseText.indexOf("{") == 0 ){
            document.getElementById("interceptedResponse").value =
            self.responseText;
        }  
      }
      if(oldOnReadyStateChange) {
        oldOnReadyStateChange();
      }
    }

    if(this.addEventListener) {
      this.addEventListener("readystatechange", onReadyStateChange,
        false);
    } else {
      oldOnReadyStateChange = this.onreadystatechange;
      this.onreadystatechange = onReadyStateChange;
    }
    send.call(this, data);
  }
})(XMLHttpRequest);`

const (
	//设置常量 分别设置chromedriver.exe的地址和本地调用端口
	seleniumPath1 = `D:\pc\chromedriver.exe`
	seleniumPath2 = `D:\pc1\chromedriver.exe`
	port          = 9515
	SleepTime     = 5
)

var Stop = true

// StartChrome 启动谷歌浏览器headless模式
func StartChrome(url string) {

	Stop = false
	//开两个页面，头尾一起
	wd1, service1 := GetDriver(url, 9515, seleniumPath1)
	//延迟退出chrome
	go GoRoutineStart(wd1, service1)

	wd2, service2 := GetDriver(url, 9515, seleniumPath1)
	go GoRoutineEnd(wd2, service2)

	var PageStart = 0
	var PageEnd = 50

	for {
		select {
		case stop := <-Manger.Reponse:
			if stop.Type == 0 { //从开始读
				PageStart = stop.Current_page
				fmt.Println("开始读到 ", PageStart)
			} else { //从结束读
				PageEnd = stop.Current_page
				fmt.Println("末尾读到 ", PageEnd)
			}
			if PageStart >= PageEnd {
				fmt.Println("PC ----- finished")
				Stop = true
				return
			}

		}
	}

}

func GetDriver(url string, port int, path string) (selenium.WebDriver, *selenium.Service) {
	//1.开启selenium服务
	//设置selium服务的选项,设置为空。根据需要设置。
	ops := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(path, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}
	// //延迟关闭服务
	// defer service.Stop()

	//2.调用浏览器
	//设置浏览器兼容性，我们设置浏览器名称为chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	// //延迟退出chrome
	// defer wd.Quit()

	//3.对页面元素进行操作
	if err := wd.Get(url); err != nil {
		panic(err)
	}

	return wd, service
}

//从第二页开始
func GoRoutineStart(wd selenium.WebDriver, service *selenium.Service) {
	//执行js
	var args []interface{}
	args = make([]interface{}, 0)
	result, err := wd.ExecuteScript(GetPostData, args)
	//睡眠20秒后退出
	time.Sleep(5 * time.Second)
	fmt.Println(result, err)

	//先跳转下一页
	wd, _, err = ClickToNextPage(wd)
	if err != nil {
		panic(err)
	}

	//再跳转第一页,拿第一页的数据
	wd, PageInfo, err := ClickToUpPage(wd)

	if err != nil {
		panic(err)
	}

	Manger.Request <- RequestCh{Data: PageInfo, Type: 0}

	for {
		wd, PageInfo, err = ClickToNextPage(wd)
		Manger.Request <- RequestCh{Data: PageInfo, Type: 0}
		if Stop {
			defer wd.Quit()
			defer service.Stop()
			return
		}
	}

}

//从最后一页开始
func GoRoutineEnd(wd selenium.WebDriver, service *selenium.Service) {

	//执行js
	var args []interface{}
	args = make([]interface{}, 0)
	result, err := wd.ExecuteScript(GetPostData, args)
	//睡眠20秒后退出
	time.Sleep(5 * time.Second)
	fmt.Println(result, err)

	//先跳转最后一页
	wd, PageInfo, err := ClickToLastPage(wd)
	if err != nil {
		panic(err)
	}

	Manger.Request <- RequestCh{Data: PageInfo, Type: 1}

	for {
		wd, PageInfo, err = ClickToUpPage(wd)
		Manger.Request <- RequestCh{Data: PageInfo, Type: 1}
		if Stop {
			defer wd.Quit()
			defer service.Stop()
			return
		}
	}

}

//点击最后一页
func ClickToLastPage(w selenium.WebDriver) (selenium.WebDriver, []byte, error) {

	var inbyte []byte

	ws, err := w.FindElements(selenium.ByCSSSelector, ".page-item")

	if err != nil {
		return w, inbyte, err
	}

	if len(ws) < 0 {
		return w, inbyte, err
	}

	err = ws[len(ws)-2].Click()
	utils.RandSleep(SleepTime)

	//拿第一页input内容
	input, err := w.FindElement(selenium.ByXPATH, `//*[@id="interceptedResponse"]`)
	if err != nil {
		panic(err)
	}
	jsinput, err := input.GetAttribute("value")
	if err != nil {
		panic(err)
	}

	return w, []byte(jsinput), err
}

//点击到下一页
func ClickToNextPage(w selenium.WebDriver) (selenium.WebDriver, []byte, error) {
	var inbyte []byte

	ws, err := w.FindElements(selenium.ByCSSSelector, ".page-item")

	if err != nil {
		return w, inbyte, err
	}

	if len(ws) < 0 {
		fmt.Println("未有下一页标签")
		return w, inbyte, err
	}

	err = ws[len(ws)-1].Click()
	utils.RandSleep(SleepTime)

	//拿第一页input内容
	input, err := w.FindElement(selenium.ByXPATH, `//*[@id="interceptedResponse"]`)
	if err != nil {
		panic(err)
	}
	jsinput, err := input.GetAttribute("value")
	if err != nil {
		panic(err)
	}

	return w, []byte(jsinput), err
}

//点击到上一页
func ClickToUpPage(w selenium.WebDriver) (selenium.WebDriver, []byte, error) {

	var inbyte []byte

	ws, err := w.FindElements(selenium.ByCSSSelector, ".page-item")

	if err != nil {
		return w, inbyte, err
	}

	if len(ws) < 0 {
		fmt.Println("未有上一页标签")
		return w, inbyte, err
	}

	err = ws[0].Click()
	utils.RandSleep(SleepTime)

	//拿第一页input内容
	input, err := w.FindElement(selenium.ByXPATH, `//*[@id="interceptedResponse"]`)
	if err != nil {
		panic(err)
	}
	jsinput, err := input.GetAttribute("value")
	if err != nil {
		panic(err)
	}

	return w, []byte(jsinput), err
}
