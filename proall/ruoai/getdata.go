package ruoai

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"test/model/ruoai"
	"test/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"github.com/tebeka/selenium"
)

//城市列表解析
func ParseCity(url string, chananme int) ParseResult {
	var pre ParseResult
	wd, service := GetDriver(url, Port+chananme, SeleniumPath)
	defer service.Stop()
	defer wd.Close()
	w1, err := wd.FindElement(selenium.ByCSSSelector, ".layui-table")
	if err != nil {
		return pre
	}
	w2, err := w1.FindElements(selenium.ByTagName, "tr")

	if err != nil {
		return pre
	}
	for _, value := range w2 {

		fromstr, err := value.Text()
		if err != nil {
			continue
		}
		from := strings.Split(fromstr, " ")
		fmt.Println(from)

		w3, err := value.FindElements(selenium.ByTagName, "a")
		if err != nil || len(from) < 2 {
			continue
		}

		if len(w3) > 0 && len(w3) == len(from)-1 {

			for key, value := range w3 {

				dpro, err := ruoai.GetProvinceFromRedis(from[0]) //先查省
				if err != nil && err != gorm.ErrRecordNotFound {
					continue
				}
				if dpro["id"] == "" {
					if ruoai.ProvinceCreate(from[0]) != nil {
						continue
					}
					dpro, err = ruoai.GetProvinceFromRedis(from[0])
				}
				dcity, err := ruoai.GetCityFromRedis(from[key+1]) //再查市
				if err != nil && err != gorm.ErrRecordNotFound {
					continue
				}

				if dcity["id"] == "" {
					proid, _ := strconv.Atoi(dpro["id"])
					if ruoai.CityCreate(from[key+1], proid) != nil {
						continue
					}
				}
				url, err = value.GetAttribute("href")
				ruoai.PutCityData(url)
			}
		}

		//pre.Requets = append(pre.Requets, Request{Url: })
	}

	return pre
}

//废弃
func GetAllUserUrl(url string, chananme int) ParseResult {
	var pregirl ParseResult
	var pregboy ParseResult
	var group sync.WaitGroup

	go func() {
		group.Add(1)
		pregirl = GetBoyUrl(url, 1)
		group.Done()
	}()

	go func() {
		group.Add(1)
		pregboy = GetGirlUrl(url, 2)
		group.Done()
	}()

	group.Wait()

	for _, value := range pregirl.Requets {
		pregboy.Requets = append(pregboy.Requets, value)
	}

	return pregboy
}

func GetBoyUrl(url string, chananme int) ParseResult {
	var pre ParseResult
	wd, service := GetDriver(url, Port+chananme, SeleniumPath)
	defer service.Stop()
	defer wd.Close()
	//男按钮
	w1, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/div/div[1]/div[1]")
	if err != nil {
		return pre
	}
	err = w1.Click()
	if err != nil {
		return pre
	}
	w2, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/div/div[3]/button")
	if err != nil {
		return pre
	}
	err = w2.Click()
	if err != nil {
		return pre
	}

	w3, err := wd.FindElements(selenium.ByCSSSelector, ".list_a_avatar")
	for _, value := range w3 {
		value1, err := value.GetAttribute("href")
		if err != nil {
			continue
		}
		pre.Requets = append(pre.Requets, Request{Url: value1, ParserFunc: GetInformation})
	}

	//下一页
	w4, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[3]/div[1]/ul/ul/li[2]/a")
	if err != nil {
		return pre
	}

	nextPageurl, err := w4.GetAttribute("href")
	if err != nil {
		return pre
	}

	for {

		if nextPageurl == "" {
			break
		}

		Info, _ := utils.Fetch(nextPageurl)
		fmt.Println("原url为：" + url + "男" + nextPageurl)
		nextPageurl = ""
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(Info)))
		if err != nil {
			break
		}

		dom.Find(".wrap").Each(func(i int, t *goquery.Selection) {
			t.Find(".clearfix").Each(func(i int, t *goquery.Selection) {
				t.Find(".list_a_avatar").Each(func(i int, t *goquery.Selection) {
					href, _ := t.Attr("href")
					pre.Requets = append(pre.Requets, Request{Url: href, ParserFunc: GetInformation})
				})

			})
		})

		nextPage := []string{}
		dom.Find(".footer").Each(func(i int, t *goquery.Selection) {
			t.Find(".pager").Each(func(i int, t *goquery.Selection) {
				t.Find("a").Each(func(i int, t *goquery.Selection) {
					href, _ := t.Attr("href")
					nextPage = append(nextPage, href)
				})
			})
		})

		if len(nextPage) > 1 {
			nextPageurl = RuoaiDomain + nextPage[1]
		}

	}
	fmt.Println("返回")

	return pre
}

func GetGirlUrl(url string, chananme int) ParseResult {
	var pre ParseResult
	wd, service := GetDriver(url, Port+chananme, SeleniumPath)
	defer service.Stop()
	defer wd.Close()
	//女按钮
	w1, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/div/div[1]/div[1]")
	if err != nil {
		return pre
	}
	err = w1.Click()
	if err != nil {
		return pre
	}
	w2, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[2]/form/div/div[1]/div[2]/i")
	if err != nil {
		return pre
	}
	err = w2.Click()
	if err != nil {
		return pre
	}

	w3, err := wd.FindElements(selenium.ByCSSSelector, ".list_a_avatar")
	for _, value := range w3 {
		value1, err := value.GetAttribute("href")
		if err != nil {
			continue
		}
		pre.Requets = append(pre.Requets, Request{Url: value1, ParserFunc: GetInformation})
	}

	//下一页
	w4, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[3]/div[1]/ul/ul/li[2]/a")
	if err != nil {
		return pre
	}

	nextPageurl, err := w4.GetAttribute("href")
	if err != nil {
		return pre
	}

	for {

		if nextPageurl == "" {
			break
		}

		Info, _ := utils.Fetch(nextPageurl)
		fmt.Println("原url为：" + url + "女" + nextPageurl)
		nextPageurl = ""
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(Info)))
		if err != nil {
			break
		}

		dom.Find(".wrap").Each(func(i int, t *goquery.Selection) {
			t.Find(".clearfix").Each(func(i int, t *goquery.Selection) {
				t.Find(".list_a_avatar").Each(func(i int, t *goquery.Selection) {
					href, _ := t.Attr("href")
					pre.Requets = append(pre.Requets, Request{Url: href, ParserFunc: GetInformation})
				})

			})
		})

		nextPage := []string{}
		dom.Find(".footer").Each(func(i int, t *goquery.Selection) {
			t.Find(".pager").Each(func(i int, t *goquery.Selection) {
				t.Find("a").Each(func(i int, t *goquery.Selection) {
					href, _ := t.Attr("href")
					nextPage = append(nextPage, href)
				})
			})
		})

		if len(nextPage) > 1 {
			nextPageurl = RuoaiDomain + nextPage[1]
		}

	}
	fmt.Println("返回")
	return pre
}

func GetMoney(number, unit string) int {
	value := 0
	if unit == "千" {
		value, _ = strconv.Atoi(number)
		return value * 1000
	}
	if unit == "万" {
		value, _ = strconv.Atoi(number)
		return value * 10000
	}
	return value
}

func WriteUser(u ruoai.User) {

	dpro, err := ruoai.GetUserFromRedis(strconv.Itoa(u.Id)) //
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if dpro["id"] == "" {
		if ruoai.UserCreate(u) != nil {
			return
		}
	}
}

func GetInformation(url string, i int) ParseResult {

	Info, _ := utils.Fetch(url)
	var user ruoai.User
	var pre ParseResult
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(Info)))
	if err != nil {
		fmt.Println("ParseMerInfoList  ", err.Error())
	}

	dom.Find(".intro").Each(func(i int, selection *goquery.Selection) {
		selection.Find(".albums").Each(func(n int, t *goquery.Selection) {
			user.Url_image, _ = t.Find("img").Eq(0).Attr("src")
		})
		var Info []string
		selection.Find(".basic_info").Each(func(i int, t *goquery.Selection) {
			user.Nickname = t.Find(".nickname").Text()
			t.Find(".simple_info > span").Each(func(n int, t *goquery.Selection) {
				fmt.Println(t.Text())
				Info = append(Info, t.Text())
			})
			if len(Info) < 4 {
				return
			}
			//ID
			reid := "ID：([0-9]+)"
			re := regexp.MustCompile(reid)
			math := re.FindAllSubmatch([]byte(Info[0]), -1)
			user.Id, _ = strconv.Atoi(string(math[0][1]))
			//年龄
			reage := "([0-9]+)岁"
			re = regexp.MustCompile(reage)
			math = re.FindAllSubmatch([]byte(Info[1]), -1)
			if len(math[0]) > 1 {
				user.Age, _ = strconv.Atoi(string(math[0][1]))
			}
			//性别+学历
			user.Gender = Info[2]
			dpro, err := ruoai.GetEducationFromRedis(Info[3]) //
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println("用户页面数据异常", url)
				return
			}
			if dpro["id"] == "" {
				if ruoai.EducationCreate(Info[3]) != nil {
					fmt.Println("用户页面数据异常", url)
					return
				}
				dpro, err = ruoai.GetEducationFromRedis(Info[3])
			}
			user.Education, err = strconv.Atoi(dpro["id"])
		})

	})

	var allinfo []string
	dom.Find(".more_info_content").Each(func(i int, t *goquery.Selection) {
		t.Find(".show_info_content").Each(func(i int, t *goquery.Selection) {
			t.Find("span").Each(func(i int, t *goquery.Selection) {
				fmt.Println(t.Text())
				allinfo = append(allinfo, t.Text())
			})
		})
		user.Title = t.Find(".idea_content").Text()
	})

	for key, value := range allinfo {
		values := strings.Split(value, "：")
		switch key {
		case 0:
			//居住地
			if len(values) > 1 {
				livecity, _ := ruoai.GetCityFromRedis(values[1])
				user.City_code, _ = strconv.Atoi(livecity["id"])
			}
		case 1:
			//生日
			if len(values) > 1 {
				user.Birthday_time = values[1]
			}
		case 2:
			//身高
			if len(values) > 1 {
				user.Height = values[1]
			}
		case 3:
			//体重
			if len(values) > 1 {
				user.Weight = values[1]
			}
		case 4:
			//收入
			if len(values) > 1 {

				if !strings.Contains(values[1], "千") && !strings.Contains(values[1], "万") {
					break
				}

				if strings.Contains(values[1], "-") {
					reage := `([0-9]+)([^\d]) - ([0-9]+)([^\d])*`
					re := regexp.MustCompile(reage)
					math := re.FindAllSubmatch([]byte(values[1]), -1)
					user.Salary_down = GetMoney(string(math[0][1]), string(math[0][2]))
					user.Salary_up = GetMoney(string(math[0][3]), string(math[0][4]))
					break
				}

				if strings.Contains(values[1], "以上") {
					user.Salary_down = 0
					reage := `([0-9]+)([^\d])以上`
					re := regexp.MustCompile(reage)
					math := re.FindAllSubmatch([]byte(values[1]), -1)
					user.Salary_up = GetMoney(string(math[0][1]), string(math[0][2]))
					break
				}

				if strings.Contains(values[1], "以下") {
					user.Salary_up = 0
					reage := `([0-9]+)([^\d])以下`
					re := regexp.MustCompile(reage)
					math := re.FindAllSubmatch([]byte(values[1]), -1)
					user.Salary_down = GetMoney(string(math[0][1]), string(math[0][2]))
					break
				}

			}
		case 5:
			//婚姻状态
			if len(values) > 1 {
				user.Married = values[1]
			}
		case 6:
			//购房
			if len(values) > 1 {
				user.House = values[1]
			}
		case 8:
			//职业
			if len(values) > 1 {
				if values[1] != "" {
					dpro, err := ruoai.GetWorkFromRedis(values[1]) //
					if err != nil && err != gorm.ErrRecordNotFound {
						fmt.Println("用户页面数据异常", url)
						return pre
					}
					if dpro["id"] == "" {
						if ruoai.WorkCreate(values[1]) != nil {
							fmt.Println("用户页面数据异常", url)
							return pre
						}
						dpro, err = ruoai.GetWorkFromRedis(values[1])
					}
					user.Work, err = strconv.Atoi(dpro["id"])
				}
			}
		case 9:
			//家乡
			if len(values) > 1 {
				livecity, _ := ruoai.GetCityFromRedis(values[1])
				user.Hometown, _ = strconv.Atoi(livecity["id"])
			}
		}
	}
	WriteUser(user)
	return pre

}

//获取页面信息
func GetInformation1(url string, chananme int) ParseResult {
	var pre ParseResult
	var user ruoai.User
	wd, service := GetDriver(url, Port+chananme, SeleniumPath)
	defer service.Stop()
	defer wd.Close()

	//头像
	w1, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[2]/div[2]/div[1]/img")
	if err != nil {
		return pre
	}
	url_imge, err := w1.GetAttribute("src")
	if err == nil {
		user.Url_image = url_imge
	}

	w2, err := wd.FindElement(selenium.ByCSSSelector, ".basic_info")
	if err != nil {
		return pre
	}

	str, err := w2.Text()
	if err != nil {
		return pre
	}
	if len(str) < 3 {
		fmt.Println("用户页面数据异常", url)
		return pre
	}
	strfull := strings.Split(str, "\n")
	if len(strfull) < 3 {
		strs := []string{}
		strs = append(strs, "")
		for _, value := range strfull {
			strs = append(strs, value)
		}
		strfull = strs
	}
	user.Nickname = strfull[0]
	id_age := strings.Split(strfull[1], " ")
	//ID
	reid := "ID：([0-9]+)"
	re := regexp.MustCompile(reid)
	math := re.FindAllSubmatch([]byte(id_age[0]), -1)
	user.Id, _ = strconv.Atoi(string(math[0][1]))

	//年龄
	reage := "([0-9]+)岁"
	re = regexp.MustCompile(reage)
	math = re.FindAllSubmatch([]byte(id_age[1]), -1)
	if len(math[0]) > 1 {
		user.Age, _ = strconv.Atoi(string(math[0][1]))
	}
	//性别+学历
	gender_sc := strings.Split(strfull[2], " ")
	user.Gender = gender_sc[0]

	if gender_sc[1] != "" {

		dpro, err := ruoai.GetEducationFromRedis(gender_sc[1]) //
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println("用户页面数据异常", url)
			return pre
		}
		if dpro["id"] == "" {
			if ruoai.EducationCreate(gender_sc[1]) != nil {
				fmt.Println("用户页面数据异常", url)
				return pre
			}
			dpro, err = ruoai.GetEducationFromRedis(gender_sc[1])
		}
		user.Education, err = strconv.Atoi(dpro["id"])
	}

	//详细信息
	w3, err := wd.FindElements(selenium.ByCSSSelector, ".show_info_content")
	if err != nil {
		return pre
	}

	if len(w3) < 2 {
		fmt.Println("用户页面数据异常", url)
		return pre
	}
	str, err = w3[0].Text()
	if err != nil {
		fmt.Println("用户页面数据异常", url)
		return pre
	}
	strfull = strings.Split(str, "\n")
	for key, value := range strfull {
		values := strings.Split(value, "：")
		switch key {
		case 0:
			//居住地
			if len(values) > 1 {
				livecity, _ := ruoai.GetCityFromRedis(values[1])
				user.City_code, _ = strconv.Atoi(livecity["id"])
			}
		case 1:
			//生日
			if len(values) > 1 {
				user.Birthday_time = values[1]
			}
		case 2:
			//身高
			if len(values) > 1 {
				user.Height = values[1]
			}
		case 3:
			//体重
			if len(values) > 1 {
				user.Weight = values[1]
			}
		case 4:
			//收入
			if len(values) > 1 {

				if !strings.Contains(values[1], "千") && !strings.Contains(values[1], "万") {
					break
				}

				if strings.Contains(values[1], "-") {
					reage := `([0-9]+)([^\d]) - ([0-9]+)([^\d])*`
					re = regexp.MustCompile(reage)
					math = re.FindAllSubmatch([]byte(values[1]), -1)
					user.Salary_down = GetMoney(string(math[0][1]), string(math[0][2]))
					user.Salary_up = GetMoney(string(math[0][3]), string(math[0][4]))
					break
				}

				if strings.Contains(values[1], "以上") {
					user.Salary_down = 0
					reage := `([0-9]+)([^\d])以上`
					re = regexp.MustCompile(reage)
					math = re.FindAllSubmatch([]byte(values[1]), -1)
					user.Salary_up = GetMoney(string(math[0][1]), string(math[0][2]))
					break
				}

				if strings.Contains(values[1], "以下") {
					user.Salary_up = 0
					reage := `([0-9]+)([^\d])以下`
					re = regexp.MustCompile(reage)
					math = re.FindAllSubmatch([]byte(values[1]), -1)
					user.Salary_down = GetMoney(string(math[0][1]), string(math[0][2]))
					break
				}

			}
		case 5:
			//婚姻状态
			if len(values) > 1 {
				user.Married = values[1]
			}
		case 6:
			//购房
			if len(values) > 1 {
				user.House = values[1]
			}
		case 8:
			//职业
			if len(values) > 1 {
				if values[1] != "" {
					dpro, err := ruoai.GetWorkFromRedis(values[1]) //
					if err != nil && err != gorm.ErrRecordNotFound {
						fmt.Println("用户页面数据异常", url)
						return pre
					}
					if dpro["id"] == "" {
						if ruoai.WorkCreate(values[1]) != nil {
							fmt.Println("用户页面数据异常", url)
							return pre
						}
						dpro, err = ruoai.GetWorkFromRedis(values[1])
					}
					user.Work, err = strconv.Atoi(dpro["id"])
				}
			}
		case 9:
			//家乡
			if len(values) > 1 {
				livecity, _ := ruoai.GetCityFromRedis(values[1])
				user.Hometown, _ = strconv.Atoi(livecity["id"])
			}
		}
	}

	//年龄
	str, err = w3[1].Text()
	if err != nil {
		fmt.Println("用户页面数据异常", url)
		return pre
	}
	strfull = strings.Split(str, "\n")
	for _, value := range strfull {

		values := strings.Split(value, "：")
		reage := `([0-9]+)岁`
		re = regexp.MustCompile(reage)
		math = re.FindAllSubmatch([]byte(values[1]), -1)
		user.Age, _ = strconv.Atoi(string(math[0][1]))
		break
	}

	//爱情宣言
	w4, err := wd.FindElement(selenium.ByCSSSelector, ".idea_content")
	if err != nil {
		return pre
	}
	str, err = w4.Text()
	if err != nil {
		return pre
	}
	user.Title = str
	fmt.Println("user 信息", user)
	WriteUser(user)
	return pre
}
