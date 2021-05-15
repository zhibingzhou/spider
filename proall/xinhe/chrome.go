package xinhe

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func GetVideo(VideoUrl string, film_type int) PareResult {

	var pre PareResult

	film_ids := strings.Split(VideoUrl, "/")
	if len(film_ids) < 2 {
		return pre
	}
	id, _ := strconv.Atoi(film_ids[len(film_ids)-1])

	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	var Stop chan string
	Stop = make(chan string)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("window-size", "50,400"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
		chromedp.UserDataDir(dir),
	)

	c, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel() //浏览器不关闭，可以打开这个
	// 给每个页面的爬取设置超时时间
	chromeCtx, cancel = context.WithTimeout(chromeCtx, 20*time.Second)
	defer cancel()

	//defer chromedp.Cancel(chromeCtx)

	// 执行一个空task, 用提前创建Chrome实例
	// ensure that the browser process is started
	if err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...); err != nil {
		panic(err)
	}

	// listen network event
	listenForNetworkEvent(chromeCtx, Stop)

	var htmls string

	go func() {

		if err := chromedp.Run(chromeCtx,
			network.Enable(),
			chromedp.Navigate(VideoUrl),
			chromedp.WaitVisible(`#__next > div.jsx-2684177236 > div > div.jsx-2803742412`, chromedp.BySearch),
			chromedp.OuterHTML(`/html/body`, &htmls, chromedp.BySearch),
		); err != nil {
			fmt.Println(err)
		}

	}()

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case url := <-Stop:
			PareFilm(htmls, url, film_type, id)
			return pre
		case <-ticker.C:
			fmt.Println("超时")
			return pre
		}

	}
}

//监听
func listenForNetworkEvent(ctx context.Context, c chan string) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {

		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {
				// log.Printf("received headers: %s", resp.Headers)
				if strings.Index(resp.URL, ".m3u8") != -1 {
					log.Printf("received headers: %s", resp.URL)
					c <- resp.URL
					return
				}
			}
			fmt.Println(resp.URL)
		}
		// other needed network Event
	})
}

//模拟登录
func GetChrome(url string) {
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

	//执行命令
	// var htmlContent string
	// err := chromedp.Run(timeoutCtx,
	// 	chromedp.Navigate(link),
	// 	chromedp.WaitVisible(waitExpression),
	// 	chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
	// )

	var buf []byte
	var buf1 []byte
	var htmls string
	if err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.BySearch),
		chromedp.SendKeys(`name`, "konwen", chromedp.ByID),
		chromedp.SendKeys(`password`, "123456", chromedp.ByID),
		chromedp.Screenshot(`document.querySelector("#root")`, &buf, chromedp.ByJSPath),
		chromedp.Click(`//*[@id="root"]/div[3]/div/div[2]/form/div[3]/div/button`, chromedp.BySearch),
		chromedp.OuterHTML(`.C-Header`, &htmls, chromedp.BySearch), //XPATH .class
		chromedp.Sleep(time.Duration(3)*time.Second),
		chromedp.WaitVisible(`body`, chromedp.BySearch),
		//	chromedp.WaitReady(`#root > div.container`, chromedp.BySearch),
		// chromedp.Screenshot(`document.querySelector("#root > div.container > div")`, &buf, chromedp.ByJSPath),
		chromedp.Screenshot(`document.querySelector("#root > div.container")`, &buf1, chromedp.ByJSPath), //js path
		saveCookies(),

		//loadCookies(),
	); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("4.png", buf, 0644); err != nil {
		fmt.Println(err)
	}

	fmt.Println(htmls)

	return

}

// 保存Cookies
func saveCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {

		// cookies的获取对应是在devTools的network面板中
		// 1. 获取cookies
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return
		}

		// 2. 序列化
		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}

		// 3. 存储到临时文件
		if err = ioutil.WriteFile("cookies.tmp", cookiesData, 0755); err != nil {
			return
		}
		return
	}
}

// 加载Cookies
func loadCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 如果cookies临时文件不存在则直接跳过
		if _, _err := os.Stat("cookies.tmp"); os.IsNotExist(_err) {
			return
		}

		// 如果存在则读取cookies的数据
		cookiesData, err := ioutil.ReadFile("cookies.tmp")
		if err != nil {
			return
		}

		// 反序列化
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			return
		}

		// 设置cookies
		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

//拿到新的cookie
func Getagain(url string) {

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 给每个页面的爬取设置超时时间
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 20*time.Second)

	defer cancel()

	// 执行一个空task, 用提前创建Chrome实例
	// ensure that the browser process is started
	if err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...); err != nil {
		panic(err)
	}

	var buf3 []byte
	if err := chromedp.Run(timeoutCtx,
		loadCookies(),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.BySearch),
		chromedp.Screenshot(`//*[@id="root"]/div[1]`, &buf3, chromedp.BySearch), //js path
	); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("6.png", buf3, 0644); err != nil {
		fmt.Println(err)
	}
}
