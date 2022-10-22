package servcie

import (
	"context"
	"cosTgBot/config"
	"cosTgBot/util"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

var cp *util.LocalCipher

func InitService() {
	var err error
	cp, err = util.InitCipher([]byte("12345678901234567890123456789012"), []byte("1234567890123456"))
	if err != nil {
		panic(err)
	}

	u, _ := url.Parse(config.Get().Cos.BucketURL)
	su, _ := url.Parse(config.Get().Cos.ServiceURL)
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	cosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Get().Cos.SecretID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: config.Get().Cos.SecretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})
	_, _, err = cosClient.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}

}
