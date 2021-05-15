package xinhe

import "testing"

func TestDownload(t *testing.T) {
	//DownLoadFile()
}

func TestWriteFile(t *testing.T) {
	var m *M3u8WebFile
	m = &M3u8WebFile{Url: "https://vod.bunediy.com/20210501/mL02M2tO/index.m3u8"}
	m.GetFile()
}

func TestGetPage(t *testing.T) {
	FindScr("https://xinghe.tv/movie?type=movie&genre=%E7%81%BE%E9%9A%BE&keyword=&region=all&sort=hot&year=all", 1)
}
