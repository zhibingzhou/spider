package sezhan

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"test/model"
)

//下载路径配置
var (
	NewPath = "./newpath/"
)

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

	video_url := strings.Split(m.Url, "index.m3u8")
	m.VideoUrl = video_url[0]

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
