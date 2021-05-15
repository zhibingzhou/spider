package zhenai

import (
	"fmt"
	"test/model/zhenai"
	"test/utils"

	"github.com/jinzhu/gorm"
	"github.com/tebeka/selenium"
)

func GetCity(url string) ParseResult {
	var pre ParseResult
	wd, service := GetDriver(Port, SeleniumPath)
	defer service.Stop()
	defer wd.Close()

	err := wd.Get(url)
	if err != nil {
		return pre
	}

	utils.RandSleep(5)

	wd1, err := wd.FindElement(selenium.ByXPATH, `/html/body/div[2]/div/div[2]/img`)
	if err != nil {
		return pre
	}
	err = wd1.Click()
	if err != nil {
		return pre
	}

	wd2, err := wd.FindElement(selenium.ByXPATH, `//*[@id="app"]/div[2]/div/div[2]`)
	if err != nil {
		return pre
	}
	err = wd2.Click()
	if err != nil {
		return pre
	}

	wd3, err := wd.FindElement(selenium.ByXPATH, `//*[@id="app"]/article[2]/dl`)
	if err != nil {
		return pre
	}

	fmt.Println(wd1.Text())
	wd4, err := wd3.FindElements(selenium.ByTagName, "dd")
	if err != nil {
		return pre
	}

	var FirstW []string
	wd5, err := wd3.FindElements(selenium.ByTagName, "dt")
	if err != nil {
		return pre
	}

	for _, value := range wd5 {
		word, err := value.Text()
		if err != nil {
			continue
		}
		FirstW = append(FirstW, word)
	}

	for key, value := range wd4 {

		if len(FirstW) != len(wd4) {
			fmt.Println("长度不一致")
			return pre
		}

		he, err := value.FindElements(selenium.ByTagName, "a")
		if err != nil {
			continue
		}
		fmt.Println(FirstW[key])
		for _, value1 := range he {

			tex, err := value1.Text()
			if err != nil {
				continue
			}
			fmt.Println()

			herf, err := value1.GetAttribute("href")
			if err != nil {
				continue
			}
			dcity, err := zhenai.GetCityFromRedis(tex) //再查市
			if err != nil && err != gorm.ErrRecordNotFound {
				continue
			}

			if dcity["id"] == "" {
				if zhenai.CityCreate(tex, FirstW[key]) != nil {
					continue
				}
			}
			zhenai.PutCityData(herf)
		}

	}

	return pre
}

//找到所有页的用户
func PareAllPage(url string, number int) ParseResult {

	var wdr selenium.WebDriver
	//第一页
	pre, wd2 := GetUserUrl(url, number)
	wdr = wd2
	nextUrl := ""
	for {

		//找下一页
		wd4, err := wdr.FindElement(selenium.ByCSSSelector, ".f-pager")
		if err != nil {
			return pre
		}

		wd5, err := wd4.FindElements(selenium.ByTagName, "li")
		if err != nil {
			return pre
		}

		if len(wd5) > 0 {

			w7, err := wd5[len(wd5)-1].FindElement(selenium.ByTagName, "a")
			if err != nil {
				return pre
			}

			url, err := w7.GetAttribute("href")

			if err != nil {
				return pre
			}
			if nextUrl == url {
				fmt.Println("结束")
				return pre
			}
			nextUrl = url
			//下一页
			pre1, wd2 := GetUserUrl(url, number)
			if len(pre1.Requests) > 0 {
				pre.Requests = append(pre.Requests, pre1.Requests...)
			}
			wdr = wd2

		}
	}

	return pre
}

func GetUserUrl(url string, number int) (ParseResult, selenium.WebDriver) {
	var pre ParseResult
	err := ManagerZhenAi.Driver[number].Get(url)
	if err != nil {
		return pre, ManagerZhenAi.Driver[number]
	}

	wd2, err := ManagerZhenAi.Driver[number].FindElement(selenium.ByCSSSelector, ".g-list")
	if err != nil {
		return pre, ManagerZhenAi.Driver[number]
	}

	wd3, err := wd2.FindElements(selenium.ByCSSSelector, ".list-item")

	for _, value := range wd3 {
		wd4, err := value.FindElement(selenium.ByCSSSelector, ".photo")
		if err != nil {
			continue
		}

		wd4, err = wd4.FindElement(selenium.ByTagName, "a")
		if err != nil {
			continue
		}
		url, err := wd4.GetAttribute("href")
		if err != nil {
			continue
		}
		fmt.Println(url)
		pre.Requests = append(pre.Requests, Request{Url: url, PareFunc: PareUserInformation})
	}

	return pre, ManagerZhenAi.Driver[number]
}

func PareUserInformation(url string, number int) ParseResult {
	var pre ParseResult
	driver, _ := GetDriver(Port+2, SeleniumPath)
	err := driver.Get(url)
	if err != nil {
		return pre
	}

	w1, err := ManagerZhenAi.Driver[number].FindElement(selenium.ByXPATH, `//*[@id="app"]/div[2]/div[2]/div[1]/div[1]/div[1]/div[2]/div[1]/div[3]`)
	if err != nil {

	}

	te, err := w1.Text()

	fmt.Println(te)

	return pre
}
