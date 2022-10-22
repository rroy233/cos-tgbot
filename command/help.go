package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgUploader/util"
)

func HelpCommand(update *tgbotapi.Update) {
	text := `欢迎使用
	/help 获取帮助
	/keyboard 获取快捷键盘

以下命令需要对文件回复：
	/upload 上传到/upload文件夹下
	/css 上传到/css文件夹下
	/js 上传到/js文件夹下
	/img 上传到/img文件夹下
`
	util.SendPlainText(update, text)
}
