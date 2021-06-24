package sezhan

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"test/utils"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/tebeka/selenium"
)

func GetVideo(VideoUrl string) string {

	var re_url string

	film_ids := strings.Split(VideoUrl, "/")
	if len(film_ids) < 2 {
		return re_url
	}
	id, _ := strconv.Atoi(film_ids[len(film_ids)-1])

	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		utils.GVA_LOG.Error(err, VideoUrl)
		return re_url
	}
	defer os.RemoveAll(dir)
	var Stop chan string
	Stop = make(chan string)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		//chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("ignore-certificate-errors", true),
		//chromedp.Flag("window-size", "50,400"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
		chromedp.UserDataDir(dir),
	)

	c, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel() //浏览器不关闭，可以打开这个
	// 给每个页面的爬取设置超时时间
	chromeCtx, cancel = context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	//defer chromedp.Cancel(chromeCtx)

	// 执行一个空task, 用提前创建Chrome实例
	// ensure that the browser process is started
	if err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...); err != nil {
		utils.GVA_LOG.Error(err, VideoUrl)
		return re_url
	}

	// listen network event
	listenForNetworkEvent(chromeCtx, Stop)

	var htmls string
	var buf []byte
	go func() {

		if err := chromedp.Run(chromeCtx,
			network.Enable(),
			chromedp.Navigate(VideoUrl),
			chromedp.WaitVisible(`//*[@id="video-player-container"]/div/div[1]/div[2]`, chromedp.BySearch),
			chromedp.Click(`//*[@id="video-player-container"]/div/div[1]/div[2]`, chromedp.BySearch),//模拟点击 视频
			chromedp.Screenshot(`//*[@id="video-player-container"]/div/div[1]/div[2]`, &buf, chromedp.BySearch),
		); err != nil {
			utils.GVA_LOG.Error(err, VideoUrl)
		}

		if err := ioutil.WriteFile("4.png", buf, 0644); err != nil {
			fmt.Println(err)
		}

		fmt.Println(htmls)

	}()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case url := <-Stop: //获取链接
			utils.GVA_LOG.Debug(url, id)
			return url
		case <-ticker.C:
			fmt.Println("超时")
			utils.GVA_LOG.Debug("超时", id)
			return re_url
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

	// // 禁止加载图片，加快渲染速度
	// imagCaps := map[string]interface{}{
	// 	"profile.managed_default_content_settings.images": 2,
	// }

	// chromeCaps := chrome.Capabilities{
	// 	Prefs: imagCaps,
	// 	Path:  "",
	// 	Args: []string{
	// 		"--headless", // 设置Chrome无头模式
	// 		"--no-sandbox",
	// 		"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
	// 	},
	// }
	// caps.AddChrome(chromeCaps)

	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}

	// // 这是目标网站留下的坑，不加这个在linux系统中会显示手机网页，每个网站的策略不一样，需要区别处理。
	// webDriver.AddCookie(&selenium.Cookie{
	// 	Name:  "defaultJumpDomain",
	// 	Value: "www",
	// })

	// //延迟退出chrome
	// defer wd.Quit()

	//3.对页面元素进行操作
	if err := wd.Get(url); err != nil {
		panic(err)
	}

	return wd, service
}

func Test() {

}
