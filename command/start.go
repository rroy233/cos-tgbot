package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgUploader/util"
)

func StartCommand(update *tgbotapi.Update) {
	text := `输入 /help 获取帮助
`
	util.SendPlainText(update, text)
}
