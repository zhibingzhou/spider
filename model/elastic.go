package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

var EsRuoai EsRuoAi

type EsRuoAi struct {
	ElaClient *elastic.Client
	Index     string
}

func init() {
	NewElastic()
}

func NewElastic() {
	// 创建client
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"),
		// 基于http base auth验证机制的账号和密码
		elastic.SetBasicAuth("user", "secret"),
		// 启用gzip压缩
		elastic.SetGzip(true),
		// 设置监控检查时间间隔
		elastic.SetHealthcheckInterval(10*time.Second),
		// 设置请求失败最大重试次数
		elastic.SetMaxRetries(5),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)

	if err != nil {
		// Handle error
		fmt.Printf("连接失败: %v\n", err)
	} else {
		fmt.Println("连接成功")
	}
	EsRuoai = EsRuoAi{ElaClient: client, Index: "user"}
}
