package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigStruct struct {
	BotToken string  `json:"bot_token"`
	AdminUid []int64 `json:"admin_uid"`
	Cos      struct {
		BucketURL    string `json:"BucketURL"`
		ServiceURL   string `json:"ServiceURL"`
		SecretID     string `json:"SecretID"`
		SecretKey    string `json:"SecretKey"`
		CdnUrlDomain string `json:"cdnUrlDomain"`
	} `json:"cos"`
	Logger struct {
		Enabled   bool   `json:"enabled"`
		Report    bool   `json:"report"`
		ReportUrl string `json:"reportUrl"`
		QueryKey  string `json:"queryKey"`
	} `json:"logger"`
}

var config *ConfigStruct
var adminMap map[int64]int64

func LoadConfig() {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalln(err)
	}
	config = new(ConfigStruct)
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln(err)
	}

	//admin-map
	adminMap = make(map[int64]int64)
	for _, uid := range config.AdminUid {
		adminMap[uid] = 1
	}
}

func Get() *ConfigStruct {
	return config
}

func IsAdmin(uid int64) bool {
	if adminMap[uid] == 1 {
		return true
	}
	return false
}
