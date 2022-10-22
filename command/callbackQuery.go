package command

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strings"
	"tgUploader/servcie"
	"tgUploader/util"
)

func DelQuery(update *tgbotapi.Update) {
	text := update.CallbackQuery.Message.Text
	objKey, err := util.MatchSingle(regexp.MustCompile(`ObjectKey:【(.+?)】`), text)
	Key, err := util.MatchSingle(regexp.MustCompile(`\nKey:【(.+?)】\n`), text)
	log.Printf("objKey:%s   Key:%s", objKey, Key)
	if err != nil {
		if err := util.CallBackWithAlert(update.CallbackQuery.ID, "objKey和Sign解析失败"); err != nil {
			log.Println("[command][callbackQuery.DelQuery]objKey和Sign解析失败发送失败", err)
		}
		return
	}

	data := strings.Split(Key, "#")
	if len(data) != 2 {
		if err := util.CallBackWithAlert(update.CallbackQuery.ID, "请求参数无效"); err != nil {
			log.Println("[command][callbackQuery.DelQuery]default发送CallBackWithAlert失败", err)
		}
		return
	}

	sign := util.Sha256Hex([]byte(fmt.Sprintf("del#%s#%s", objKey, data[0])))[:6]
	if data[1] != sign {
		if err := util.CallBackWithAlert(update.CallbackQuery.ID, "请求参数签名无效"); err != nil {
			log.Println("[command][callbackQuery.DelQuery]default发送CallBackWithAlert失败", err)
		}
		return
	}

	if servcie.CosFileExist(objKey) == true {
		err := servcie.CosFileDel(objKey)
		if err != nil {
			if err := util.CallBackWithAlert(update.CallbackQuery.ID, "删除失败"); err != nil {
				log.Println("[command][callbackQuery.DelQuery]default发送CallBackWithAlert失败", err)
			}
			return
		}
	}

	if err := util.CallBack(update.CallbackQuery.ID, "成功"); err != nil {
		log.Println("[command][callbackQuery.DelQuery]default发送CallBack失败", err)
	}

	util.EditMessageText(update.CallbackQuery.Message, "文件已删除")

	return
}
