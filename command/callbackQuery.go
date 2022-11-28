package command

import (
	"cosTgBot/servcie"
	"cosTgBot/util"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/logger"
	"regexp"
	"strings"
)

func DelQuery(update *tgbotapi.Update) {
	text := update.CallbackQuery.Message.Text
	objKey, err := util.MatchSingle(regexp.MustCompile(`ObjectKey:【(.+?)】`), text)
	Key, err := util.MatchSingle(regexp.MustCompile(`\nKey:【(.+?)】\n`), text)
	logger.Info.Printf("objKey:%s   Key:%s", objKey, Key)
	if err != nil {
		if err := util.CallBackWithAlert(update.CallbackQuery.ID, "objKey和Sign解析失败"); err != nil {
			logger.Info.Println("[command][callbackQuery.DelQuery]objKey和Sign解析失败发送失败", err)
		}
		return
	}

	data := strings.Split(Key, "#")
	if len(data) != 2 {
		if err := util.CallBackWithAlert(update.CallbackQuery.ID, "请求参数无效"); err != nil {
			logger.Info.Println("[command][callbackQuery.DelQuery]default发送CallBackWithAlert失败", err)
		}
		return
	}

	sign := util.Sha256Hex([]byte(fmt.Sprintf("del#%s#%s", objKey, data[0])))[:6]
	if data[1] != sign {
		if err := util.CallBackWithAlert(update.CallbackQuery.ID, "请求参数签名无效"); err != nil {
			logger.Info.Println("[command][callbackQuery.DelQuery]default发送CallBackWithAlert失败", err)
		}
		return
	}

	if servcie.CosFileExist(objKey) == true {
		err := servcie.CosFileDel(objKey)
		if err != nil {
			if err := util.CallBackWithAlert(update.CallbackQuery.ID, "删除失败"); err != nil {
				logger.Info.Println("[command][callbackQuery.DelQuery]default发送CallBackWithAlert失败", err)
			}
			return
		}
	}

	if err := util.CallBack(update.CallbackQuery.ID, "成功"); err != nil {
		logger.Info.Println("[command][callbackQuery.DelQuery]default发送CallBack失败", err)
	}

	util.EditMessageText(update.CallbackQuery.Message, "文件已删除")

	return
}

func uploadQuery(update *tgbotapi.Update, path string) {
	
	if err := util.CallBack(update.CallbackQuery.ID, "成功"); err != nil {
		logger.Info.Println("[command][callbackQuery.uploadQuery]default发送CallBack失败", err)
	}

	UploadCommand(update.CallbackQuery.Message.ReplyToMessage, path)

	return
}

func replyToFiles(update *tgbotapi.Update) {
	sMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "正在处理...")
	sMsg.ReplyToMessageID = update.Message.MessageID
	msg, err := bot.Send(sMsg)
	if err != nil {
		logger.Info.Println("[util][callbackQuery.replyToFiles]发送正在处理信息失败:", err)
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("上传到Upload", "UPLOAD_TO_UPLOAD"),
			tgbotapi.NewInlineKeyboardButtonData("上传到Img", "UPLOAD_TO_IMG"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("上传到Css", "UPLOAD_TO_CSS"),
			tgbotapi.NewInlineKeyboardButtonData("上传到Js", "UPLOAD_TO_JS"),
		),
	)

	err = util.EditMessageWithMarkUP(&msg, "请选择上传至哪个目录喵~", &keyboard)
	if err != nil {
		logger.Info.Println("[util][callbackQuery.replyToFiles]EditMessageWithMarkUP:", err)
		util.EditMessageText(&msg, "异常")
		return
	}
}
