package sezhan

import "fmt"

/*
网站：https://jcxx77.com/video/show
网站：https://v88av.com/watch/TF9ZREw= 需要模拟点击视频
*/

type SeZhan struct {
}

type M3u8WebFile struct {
	Url      string
	VideoUrl string
	FileName string
}

func (x SeZhan) Boot() {

	url := GetVideo("https://v88av.com/watch/TF9ZREw=")

	fmt.Println(url)

	//  var m *M3u8WebFile
	//  m = &M3u8WebFile{Url: url}
	//  m.GetFile()

}
