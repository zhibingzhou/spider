package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"test/model"
	"test/proall"
	"test/router"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

)

func main() {
	//model.Delcash()
	go proall.NewProcessor().Register().Boot()
	router.Router.Run(":8082")
}

// func main() {

// 	// ctx, _ := chromedp.NewExecAllocator(
// 	// 	context.Background(),

// 	// 	// 以默认配置的数组为基础，覆写headless参数
// 	// 	// 当然也可以根据自己的需要进行修改，这个flag是浏览器的设置
// 	// 	append(
// 	// 		chromedp.DefaultExecAllocatorOptions[:],
// 	// 		chromedp.Flag("headless", false),
// 	// 		chromedp.Flag("blink-settings", "imagesEnabled=false"),
// 	// 		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
// 	// 	)...,
// 	// )

// 	// timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
// 	// defer cancel()
// 	// ctx, _ = chromedp.NewContext(
// 	// 	timeoutCtx,
// 	// 	// 设置日志方法
// 	// 	chromedp.WithLogf(log.Printf),
// 	// )

// 	// // 执行我们自定义的任务 - myTasks函数在第4步
// 	// if err := chromedp.Run(ctx, myTasks()); err != nil {
// 	// 	log.Fatal(err)
// 	// 	return
// 	// }

// 	// fmt.Println(htmlContent)

// 	//GetHttpHtmlContent("http://news.iciba.com/", "body > div.screen > div.banner > div.swiper-container-place > div > div.swiper-slide.swiper-slide-0.swiper-slide-visible.swiper-slide-active > a.item.item-big > img", `document.querySelector("body")`)
// }

var htmlContent string

// 自定义任务
func myTasks() chromedp.Tasks {
	return chromedp.Tasks{

		chromedp.ActionFunc(func(ctx context.Context) error {
			page.AddScriptToEvaluateOnNewDocument(`Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`)
			return nil
		}),

		// 1. 打开金山文档的登陆界面
		chromedp.Navigate("https://account.wps.cn/"),
		chromedp.WaitVisible("#app > div:nth-child(2) > div.CONTAINER.f-cl.f-topIndex.primary > div.CONTAINER.f-fl > div.m-userInfo > div.top.f-cl > div.logo.f-fl"),
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),

		// // 2. 点击微信登陆按钮
		// // #wechat > span:nth-child(2)
		// chromedp.Click(`#wechat > span:nth-child(2)`),

		// // 3. 点击确认按钮
		// // #dialog > div.dialog-wrapper > div > div.dialog-footer > div.dialog-footer-ok
		// chromedp.Click(`#dialog > div.dialog-wrapper > div > div.dialog-footer > div.dialog-footer-ok`),

		// 0. 加载cookies <-- 变动
		//loadCookies(),

		// 1. 打开金山文档的登陆界面

		// 判断一下是否已经登陆  <-- 变动
		//checkLoginStatus(),

		// 2. 点击微信登陆按钮
		// #wechat > span:nth-child(2)

		// 3. 点击确认按钮
		// #dialog > div.dialog-wrapper > div > div.dialog-footer > div.dialog-footer-ok

		// 4. 获取二维码
		// #wximport

		//getCode()

		// 5. 若二维码登录后，浏览器会自动跳转到用户信息页面  <-- 变动
		//saveCookies(),
	}
}

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

//获取二维码
func getCode() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 1. 用于存储图片的字节切片
		var code []byte

		// 2. 截图
		// 注意这里需要注明直接使用ID选择器来获取元素（chromedp.ByID）
		if err = chromedp.Screenshot(`#wximport`, &code, chromedp.ByID).Do(ctx); err != nil {
			return
		}

		// 3. 保存文件
		if err = ioutil.WriteFile("code.png", code, 0755); err != nil {
			return
		}

		// 3. 把二维码输出到标准输出流
		if err = printQRCode(code); err != nil {
			return err
		}

		return
	}
}

//打印二维码到终端
func printQRCode(code []byte) (err error) {
	// 1. 因为我们的字节流是图像，所以我们需要先解码字节流
	// img, _, err := image.Decode(bytes.NewReader(code))
	// if err != nil {
	// 	return
	// }

	// // 2. 然后使用gozxing库解码图片获取二进制位图
	// bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	// if err != nil {
	// 	return
	// }

	// // 3. 用二进制位图解码获取gozxing的二维码对象
	// res, err := qrcode.NewQRCodeReader().Decode(bmp, nil)
	// if err != nil {
	// 	return
	// }
	return
}

// 保存Cookies
func saveCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 等待二维码登陆
		if err = chromedp.WaitVisible(`#app`, chromedp.ByID).Do(ctx); err != nil {
			return
		}

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

// 检查是否登陆 断是否已经个人中心页面
func checkLoginStatus() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		var url string
		if err = chromedp.Evaluate(`window.location.href`, &url).Do(ctx); err != nil {
			return
		}
		if strings.Contains(url, "https://account.wps.cn/usercenter/apps") {
			log.Println("已经使用cookies登陆")
			chromedp.Stop()
		}
		return
	}
}
