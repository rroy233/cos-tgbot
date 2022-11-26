package servcie

import (
	"context"
	"fmt"
	"github.com/rroy233/logger"
	"github.com/tencentyun/cos-go-sdk-v5"
	"os"
	"time"
)

const Timeout = 5

var cosClient *cos.Client

func CosFileExist(objKey string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	resp, err := cosClient.Object.Head(ctx, objKey, nil)
	if err != nil {
		logger.Info.Println("[COS]获取文件元数据失败:", err)
		return false
	}
	contentLength := resp.Header.Get("Content-Length")

	if contentLength == "0" {
		logger.Info.Println("[COS]文件不存在" + objKey)
		return false
	}
	return true
}

func CosFileDel(objKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	_, err := cosClient.Object.Delete(ctx, objKey)
	return err
}

func CosUpload(localAddr string, path string, fileName string) (ObjectKey string, err error) {
	_, err = os.Open(localAddr)
	if err != nil {
		logger.Info.Println("[COS]尝试上传本地文件时，打开本地文件失败:", err)
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	upRes, _, err := cosClient.Object.Upload(ctx, fmt.Sprintf("%s%s", path, fileName), localAddr, nil)
	if err != nil || upRes == nil {
		logger.Info.Println("[COS]尝试上传本地文件时失败:", err)
		return "", err
	}
	logger.Info.Printf("[COS]上传文件成功:%s->%s\n", localAddr, upRes.Key)
	return upRes.Key, err
}
