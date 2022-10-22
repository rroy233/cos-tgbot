package servcie

import "tgUploader/config"

func GetFileCdnUrl(fileKey string) string {
	if config.Get().Cos.CdnUrlDomain != "" {
		return config.Get().Cos.CdnUrlDomain + fileKey
	} else {
		return config.Get().Cos.BucketURL + fileKey
	}

}
