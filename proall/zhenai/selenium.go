package zhenai

import (
	"fmt"

	//"github.com/chromedp/cdproto/page"
	"github.com/tebeka/selenium"
)

func GetDriver(port int, path string) (selenium.WebDriver, *selenium.Service) {
	//1.开启selenium服务
	//设置selium服务的选项,设置为空。根据需要设置。

	//page.AddScriptToEvaluateOnNewDocument(`Object.defineProperty(navigator, 'webdriver', {get: () => false})`)

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

	// 这是目标网站留下的坑，不加这个在linux系统中会显示手机网页，每个网站的策略不一样，需要区别处理。
	wd.AddCookie(&selenium.Cookie{
		Name:  "defaultJumpDomain",
		Value: "www",
	})

	// var args []interface{}
	// args = make([]interface{}, 0)
	// result, err := wd.ExecuteScript("Object.defineProperty(navigator, 'webdriver', { get: () => undefined })", args)
	// time.Sleep(5 * time.Second)
	// fmt.Println(result, err)

	// //延迟退出chrome
	// defer wd.Quit()

	return wd, service
}
