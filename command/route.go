package command

import (
	"cosTgBot/config"
	"cosTgBot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var bot *tgbotapi.BotAPI

func InitCommand(b *tgbotapi.BotAPI) {
	bot = b
}

func Handle(update tgbotapi.Update) {
	//判断是否包含指令
	if update.Message != nil && update.Message.IsCommand() == true {
		if config.IsAdmin(update.Message.From.ID) == false {
			util.SendPlainText(&update, "您不是管理员")
			return
		}
		//go on
		command := update.Message.Command()

		if command != "" {
			switch command {
			case "start":
				StartCommand(&update)
			case "help":
				HelpCommand(&update)
			case "upload":
				UploadCommand(&update, "/upload/")
			case "css":
				UploadCommand(&update, "/css/")
			case "js":
				UploadCommand(&update, "/js/")
			case "img":
				UploadCommand(&update, "/img/")
			default:
				util.SendPlainText(&update, "命令不存在")
				return
			}
		}

	}

	if update.CallbackQuery != nil { //判断是否为callback_query
		if config.IsAdmin(update.CallbackQuery.From.ID) == false {
			util.SendPlainText(&update, "您不是管理员")
			return
		}
		//go on
		data := update.CallbackQuery.Data
		if data == "" {
			if err := util.CallBackWithAlert(update.CallbackQuery.ID, "请求无效"); err != nil {
				log.Println("[command][route.Handle]default发送CallBackWithAlert失败", err)
			}
			return
		}
		switch update.CallbackQuery.Data {
		case "del":
			DelQuery(&update)
		default:
			if err := util.CallBackWithAlert(update.CallbackQuery.ID, "操作不存在"); err != nil {
				log.Println("[command][route.Handle]default发送CallBackWithAlert失败", err)
			}
			return
		}
	}

	//忽略
	return
}
