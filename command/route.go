package command

import (
	"cosTgBot/config"
	"cosTgBot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/logger"
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
				UploadCommand(update.Message, "/upload/")
			case "css":
				UploadCommand(update.Message, "/css/")
			case "js":
				UploadCommand(update.Message, "/js/")
			case "img":
				UploadCommand(update.Message, "/img/")
			default:
				util.SendPlainText(&update, "命令不存在")
				return
			}
		}

	}

	//是否为文件
	if update.Message != nil {
		if update.Message.Document != nil || len(update.Message.Photo) != 0 {
			//发的文件或者图片
			replyToFiles(&update)
			return
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
				logger.Info.Println("[command][route.Handle]default发送CallBackWithAlert失败", err)
			}
			return
		}
		switch update.CallbackQuery.Data {
		case "del":
			DelQuery(&update)
		case "UPLOAD_TO_UPLOAD":
			uploadQuery(&update, "/upload/")
		case "UPLOAD_TO_IMG":
			uploadQuery(&update, "/img/")
		case "UPLOAD_TO_JS":
			uploadQuery(&update, "/js/")
		case "UPLOAD_TO_CSS":
			uploadQuery(&update, "/css/")
		default:
			if err := util.CallBackWithAlert(update.CallbackQuery.ID, "操作不存在"); err != nil {
				logger.Info.Println("[command][route.Handle]default发送CallBackWithAlert失败", err)
			}
			return
		}
	}

	//忽略
	return
}
