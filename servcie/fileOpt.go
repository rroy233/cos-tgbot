package servcie

import "cosTgBot/config"

func GetFileUrl(fileKey string) string {
	return config.Get().Cos.BucketURL + fileKey
}

func GetFileCdnUrl(fileKey string) string {
	if config.Get().Cos.CdnUrlDomain != "" {
		return config.Get().Cos.CdnUrlDomain + fileKey
	} else {
		return config.Get().Cos.BucketURL + fileKey
	}
}
