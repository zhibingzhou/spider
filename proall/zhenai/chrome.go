package zhenai

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文，超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		fmt.Println("Run err : %v\n", err)
		return "", err
	}
	log.Println(htmlContent)

	return htmlContent, nil
}

func Login(url string) {

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, _ := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	//defer cancel()

	// 给每个页面的爬取设置超时时间
	chromeCtx, _ = context.WithTimeout(chromeCtx, 20*time.Second)

	//defer cancels()

	defer chromedp.Cancel(chromeCtx)

	// 执行一个空task, 用提前创建Chrome实例
	// ensure that the browser process is started
	if err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...); err != nil {
		panic(err)
	}

	var buf []byte
	//var buf1 []byte
	var htmls string
	if err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`//*[@id="login"]`, chromedp.BySearch),
		chromedp.SendKeys(`//*[@id="login"]/div/div/div[2]/section/div[1]/div[1]/input`, "17802141903", chromedp.BySearch),
		chromedp.SendKeys(`//*[@id="login"]/div/div/div[2]/section/div[2]/div/input`, "hei123xd", chromedp.BySearch),
		chromedp.Screenshot(`//*[@id="login"]`, &buf, chromedp.BySearch),
	    chromedp.Click(`//*[@id="login"]/div/div/div[2]/div/div[1]`, chromedp.BySearch),
	    chromedp.Sleep(time.Duration(10)*time.Second),
		chromedp.Screenshot(`//*[@id="login"]`, &buf, chromedp.BySearch),
		// chromedp.OuterHTML(`.C-Header`, &htmls, chromedp.BySearch), //XPATH .class
		// chromedp.Sleep(time.Duration(3)*time.Second),
		// chromedp.WaitVisible(`body`, chromedp.BySearch),
		// //	chromedp.WaitReady(`#root > div.container`, chromedp.BySearch),
		// // chromedp.Screenshot(`document.querySelector("#root > div.container > div")`, &buf, chromedp.ByJSPath),
		// chromedp.Screenshot(`document.querySelector("#root > div.container")`, &buf1, chromedp.ByJSPath), //js path
		// //saveCookies(),

		// //loadCookies(),
	); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("4.png", buf, 0644); err != nil {
		fmt.Println(err)
	}

	fmt.Println(htmls)

	return

}
