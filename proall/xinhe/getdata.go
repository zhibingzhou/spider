package xinhe

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"test/model"
	"test/model/xinhe"
	"test/utils"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type DownLoadWebFile interface {
	GetFile() string
}

func GetMovieStart(url string) {
	Info, _ := utils.Fetch(url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(Info)))
	if err != nil {
		return
	}

	dom.Find(".jsx-2684177236 > .jsx-1712581126 >.jsx-478404003> .jsx-1286128993").Each(func(i int, t *goquery.Selection) {
		isyear := false
		istype := false
		t.Find(".jsx-1913849877").Each(func(i int, t *goquery.Selection) {
			if t.Text() == "年份" {
				isyear = true
			}
			if t.Text() == "类型" {
				istype = true
			}
		})
		if isyear {
			t.Find("a").Each(func(i int, t *goquery.Selection) {

				typename := t.Text()
				fmt.Println(typename)

				herf, _ := t.Attr("href")
				fmt.Println(herf)

				if typename != "全部" {
					xinhe.PuttypeData(herf)
				}
			})
		}
		if istype {
			t.Find("a").Each(func(i int, t *goquery.Selection) {

				typename := t.Text()
				fmt.Println(typename)

				herf, _ := t.Attr("href")
				fmt.Println(herf)

				if typename != "全部" {
					filmType, _ := xinhe.GetFilmTypeFromRedis(typename)
					if filmType["id"] == "" {
						xinhe.TypeCreate(typename)
					}
				}

			})
		}
		istype = false
		isyear = false
	})
}

//一直下滑找到page面所有视频地址
func FindScr(url string, video_type int) PareResult {

	var pre PareResult
	fmt.Println(url)
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 给每个页面的爬取设置超时时间
	chromeCtx, cancel = context.WithTimeout(chromeCtx, 20*time.Second)

	defer cancel()

	//defer chromedp.Cancel(chromeCtx)

	// 执行一个空task, 用提前创建Chrome实例
	// ensure that the browser process is started
	if err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...); err != nil {
		utils.GVA_LOG.Error(err, url)
		return pre
	}

	var buf string
	var executed *runtime.RemoteObject
	if err := chromedp.Run(chromeCtx,
		chromedp.Tasks{
			chromedp.ActionFunc(func(ctxt context.Context) error {
				var html string
				chromedp.Navigate(url).Do(ctxt)
				chromedp.WaitVisible(`/html/body`, chromedp.BySearch).Do(ctxt)
				for {
					chromedp.Evaluate(`window.scrollTo(0,100000000);`, &executed).Do(ctxt) //下滑
					chromedp.Sleep(time.Duration(3 * time.Second)).Do(ctxt)
					chromedp.OuterHTML(`.infinite-scroll-component__outerdiv`, &buf, chromedp.BySearch).Do(ctxt) //js path
					if html == buf {
						break
					}
					html = buf
				}
				return nil
			}),
		},
	); err != nil {
		log.Fatal(err)
		utils.GVA_LOG.Error(err)
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(buf))
	if err != nil {
		return pre
	}

	dom.Find(".jsx-631078086").Each(func(i int, t *goquery.Selection) {
		t.Find("a").Each(func(i int, t *goquery.Selection) {
			var Info []string

			href, _ := t.Attr("href")
			fmt.Println(href)
			Info = append(Info, href)
			if href != "" {

				t.Find("img").Each(func(i int, t *goquery.Selection) {
					src, _ := t.Attr("src")
					Info = append(Info, src)
				})

				t.Find(".overlay-content-line").Each(func(i int, t *goquery.Selection) {
					fmt.Println(t.Text())
					Info = append(Info, t.Text())
				})
				Info = append(Info, t.Find(".jsx-4212652396").Text())
				CreateFilm(Info, video_type, href)

				pre.Requests = append(pre.Requests, Request{Url: href, PareFunc: GetVideo, Type: video_type})
			}
		})
	})
	utils.GVA_LOG.Debug(url, "拿到影片数：", len(pre.Requests))
	return pre
}

func CreateFilm(Info []string, video_type int, herf string) {

	if len(Info) < 9 {

		if len(Info) != 7 {
			fmt.Println("数据不全")
			utils.GVA_LOG.Debug("数据不全", herf)
			return 
		}
		var film xinhe.Film
		var film_type xinhe.Film_type
		video_id := strings.Split(Info[0], "/")
		if len(video_id) > 2 {
			film.Id, _ = strconv.Atoi(video_id[len(video_id)-1])
		}

		film.Url_image = Info[1]
		country := strings.Split(utils.DeleteExtraSpace(Info[2]), "/")
		var str_country string
		for _, value := range country {
			if value == "" {
				continue
			}
			coun := strings.ReplaceAll(value, " ", "")
			cou, _ := xinhe.GetCountryFromRedis(coun)
			if cou["id"] == "" {
				err := xinhe.CountryCreate(coun)
				if err != nil {
					fmt.Println("CreateFilm,Error", err)
					continue
				}
				cou, _ = xinhe.GetCountryFromRedis(coun)
			}
			str_country += cou["id"] + ","

		}

		film.Country = str_country

		film.First_name = Info[4]

		film.Year, _ = strconv.Atoi(Info[5])
		film.En_name = Info[6]
		film.Video_type = video_type
		film_types := strings.Split(utils.DeleteExtraSpace(Info[3]), "/")
		films, _ := xinhe.GetFilmByIdRedis(video_id[len(video_id)-1])
		if films["id"] == "" {
			err := xinhe.FilmCreate(film)
			if err != nil {
				fmt.Println("CreateFilm,Err", err)
			}
		}

		for _, value := range film_types {
			rtype, _ := xinhe.GetFilmTypeFromRedis(value)
			film_type.Type_id, _ = strconv.Atoi(rtype["id"])
			film_type.Film_id = film.Id
			err := xinhe.FilmTypeCreate(film_type.Film_id, film_type.Type_id)
			if err != nil {
				fmt.Println("CreateFilm,Err", err)
			}
		}

		return
	}

	var film xinhe.Film
	var film_type xinhe.Film_type
	video_id := strings.Split(Info[0], "/")
	if len(video_id) > 2 {
		film.Id, _ = strconv.Atoi(video_id[len(video_id)-1])
	}

	film.Url_image = Info[1]
	country := strings.Split(utils.DeleteExtraSpace(Info[2]), "/")
	var str_country string
	for _, value := range country {
		if value == "" {
			continue
		}
		coun := strings.ReplaceAll(value, " ", "")
		cou, _ := xinhe.GetCountryFromRedis(coun)
		if cou["id"] == "" {
			err := xinhe.CountryCreate(coun)
			if err != nil {
				fmt.Println("CreateFilm,Error", err)
				continue
			}
			cou, _ = xinhe.GetCountryFromRedis(coun)
		}
		str_country += cou["id"] + ","

	}

	film.Country = str_country

	film.First_name = Info[4]

	score := strings.Split(Info[5], "评分：")
	if len(score) > 1 {
		film.Score, _ = strconv.ParseFloat(score[1], 64)
	}

	film.Year, _ = strconv.Atoi(Info[7])
	film.En_name = Info[8]
	film.Video_type = video_type
	film_types := strings.Split(utils.DeleteExtraSpace(Info[3]), "/")
	films, _ := xinhe.GetFilmByIdRedis(video_id[len(video_id)-1])
	if films["id"] == "" {
		err := xinhe.FilmCreate(film)
		if err != nil {
			fmt.Println("CreateFilm,Err", err)
		}
	}

	for _, value := range film_types {
		rtype, _ := xinhe.GetFilmTypeFromRedis(value)
		film_type.Type_id, _ = strconv.Atoi(rtype["id"])
		film_type.Film_id = film.Id
		err := xinhe.FilmTypeCreate(film_type.Film_id, film_type.Type_id)
		if err != nil {
			fmt.Println("CreateFilm,Err", err)
		}
	}

}

//解析电影详细信息
func PareFilm(body string, video_url string, film_type, film_id int) {

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return
	}
	update := map[string]interface{}{}
	update["video_url"] = video_url

	dom.Find(`.jsx-526645131`).Each(func(i int, t *goquery.Selection) {
		var arrstr []string
		t.Find("a").Each(func(i int, t *goquery.Selection) {
			Starring, _ := t.Attr("title")
			arrstr = append(arrstr, Starring)
		})
		CheckInfo("标签:", arrstr, film_id)
	})

	dom.Find(`#__next > div.jsx-2684177236 > div > div.jsx-3089063832`).Each(func(i int, t *goquery.Selection) {

		t.Find(`.jsx-1308340111`).Each(func(i int, t *goquery.Selection) {
			head := t.Text()
			var arrstr []string
			t.Find("a").Each(func(i int, t *goquery.Selection) {
				Starring, _ := t.Attr("title")
				arrstr = append(arrstr, Starring)
			})
			up := CheckInfo(head, arrstr, film_id)
			for key, value := range up {
				update[key] = value
			}
		})
	})

	//更新电影信息
	xinhe.UpdateFilm(film_id, update)

}

func CheckInfo(head string, info []string, id int) map[string]interface{} {

	m := map[string]interface{}{}

	choose := strings.Split(head, ":")

	switch choose[0] {
	case "简介":
		s_name := strings.Split(head, "简介:")
		if len(s_name) > 1 {
			m["title"] = s_name[1]
		}
		break
	case "又名":
		s_name := strings.Split(head, "又名:")
		if len(s_name) > 1 {
			m["second_name"] = s_name[1]
		}
		break
	case "标签":
		for _, value := range info {
			title, _ := xinhe.GetFilmTitleFromRedis(value)
			if title["id"] == "" {
				xinhe.TitleCreate(value)
				title, _ = xinhe.GetFilmTitleFromRedis(value)
			}
			title_id, _ := strconv.Atoi(title["id"])
			filmTitle, _ := xinhe.PGetFilmTitleFromRedis(id, title_id)
			if filmTitle["id"] == "" {
				xinhe.FilmTitleCreate(id, title_id)
			}
		}
		break
	case "导演":
		for _, value := range info {
			pe, _ := xinhe.GetFilmPersonFromRedis(id, 1, value)
			if pe["id"] == "" {
				xinhe.FilmPersonCreate(id, 1, value)
			}
		}
		break
	case "编剧":
		for _, value := range info {
			pe, _ := xinhe.GetFilmPersonFromRedis(id, 2, value)
			if pe["id"] == "" {
				xinhe.FilmPersonCreate(id, 2, value)
			}
		}
		break
	case "主演":
		for _, value := range info {
			pe, _ := xinhe.GetFilmPersonFromRedis(id, 3, value)
			if pe["id"] == "" {
				xinhe.FilmPersonCreate(id, 3, value)
			}
		}
		break
	case "上映日期":
		dates := strings.Split(head, "上映日期:")
		if len(dates) > 1 {
			m["show_time"] = dates[1]
		}
		break
	}

	return m
}

//下载路径配置
var (
	NewPath = "./newpath/"
)

type M3u8WebFile struct {
	Url      string
	VideoUrl string
	FileName string
}

//通过链接 拿到 m3u8 文件
func (m *M3u8WebFile) GetFile() string {
	// Get the data
	//"https://vod.bunediy.com/20210501/mL02M2tO/index.m3u8
	resp, err := http.Get(m.Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	video_key := model.GetKey(5)
	m.FileName = NewPath + video_key + "/"

	fileInfo := strings.Split(string(body), "\n")
	var filetext = ""
	for _, line := range fileInfo {
		if strings.Contains(line, ".image") {

			fmt.Println(line)

			//下载.image文件
			DownLoadFromM3u8(m.FileName, m.VideoUrl+line)

			line = ChangeUrlHead(line) //修改链接

		} else if strings.Contains(line, ".ts") {
			fmt.Println(line)

			//下载.ts文件
			DownLoadFromM3u8(m.FileName, m.VideoUrl+line)

			line = ChangeUrlHead(line) //修改链接
		}
		filetext = filetext + line + "\n"
	}

	//重写m3u8文件
	WriteToFile(m.FileName+video_key+".m3u8", filetext)

	return m.FileName + ".m3u8"
}

//通过链接下载 .ts文件或者.image 文件
func DownLoadFromM3u8(filepath, url string) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	arrayU := strings.Split(url, "/")

	if len(arrayU) < 2 {
		return
	}

	fineName := arrayU[len(arrayU)-1]
	if ok := Exists(filepath); !ok { //创建文件夹
		err := os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			return
		}
	}

	// 创建一个文件用于保存
	out, err := os.Create(filepath + fineName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//切换路径
func ChangeUrlHead(url string) string {

	result := ""

	arrayU := strings.Split(url, "/")

	if len(arrayU) < 2 {
		return result
	}

	result = arrayU[len(arrayU)-1]

	return result
}

func WriteToFile(filePath, fileText string) {
	//创建一个新文件，写入内容 5 句 “http://c.biancheng.net/golang/”
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(fileText)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func Test() {

	url := `https://album.zhenai.com/u/1308645187`

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
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 20*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`/html/body`),
		//	chromedp.Sleep(time.Duration(1*time.Second)),
		chromedp.OuterHTML(`//*[@id="app"]`, &htmlContent, chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Run err : %v\n", err)
		return
	}
	log.Println(htmlContent)

	return

}
