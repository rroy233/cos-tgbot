package command

import (
	"cosTgBot/config"
	"cosTgBot/servcie"
	"cosTgBot/util"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/logger"
	"github.com/skip2/go-qrcode"
	"os"
	"strconv"
	"time"
)

type File struct {
	FileID   string
	FileName string
	SavePath string
	MimeType string
	FileSize int
}

const (
	CallBackQueryDel uint8 = '0'
)

func UploadCommand(message *tgbotapi.Message, path string) {
	file := new(File)
	if path == "/upload/" {
		file.SavePath = fmt.Sprintf("/upload/%d/%d/%d/%s/",
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			util.Sha256Hex([]byte(strconv.FormatInt(time.Now().UnixMilli(), 10)))[:6],
		)
	} else {
		file.SavePath = path
	}

	sMsg := tgbotapi.NewMessage(message.Chat.ID, "正在处理...")
	sMsg.ReplyToMessageID = message.MessageID
	msg, err := bot.Send(sMsg)
	if err != nil {
		logger.Info.Println("[util][upload.UploadCommand]发送正在处理信息失败:", err)
		return
	}

	if message == nil {
		util.EditMessageText(&msg, fmt.Sprintf("请给文件回复 /%s ", message.Command()))
		return
	}
	if message.ReplyToMessage != nil {
		message = message.ReplyToMessage
	}
	if message.Photo != nil && len(message.Photo) != 0 {
		file.MimeType = "image/jpeg"
		file.FileID = message.Photo[len(message.Photo)-1].FileID
		file.FileSize = message.Photo[len(message.Photo)-1].FileSize
		file.FileName = util.Sha256Hex([]byte(fmt.Sprintf("img_%d", time.Now().UnixMicro())))[:6] + ".jpg"
	} else if message.Document != nil {
		file.FileID = message.Document.FileID
		file.FileName = message.Document.FileName
		file.FileSize = message.Document.FileSize
		file.MimeType = message.Document.MimeType
	} else if message.Video != nil {
		file.FileID = message.Video.FileID
		file.FileName = util.Sha256Hex([]byte(fmt.Sprintf("img_%d", time.Now().UnixMicro())))[:6] + ".mp4"
		file.FileSize = message.Video.FileSize
		file.MimeType = message.Video.MimeType
	} else {
		util.EditMessageText(&msg, "文件不支持")
		return
	}
	logger.Info.Println("[util][upload.UploadCommand]接受文件:", util.JsonEncode(file))

	remoteFile, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: file.FileID,
	})
	if err != nil {
		logger.Info.Println("[util][upload.UploadCommand]获取文件失败:", err)
		util.EditMessageText(&msg, "获取文件失败")
		return
	}
	filePath, err := util.DownloadFile(remoteFile.Link(bot.Token))
	if err != nil {
		logger.Info.Println("[util][upload.UploadCommand]下载文件失败:", err)
		util.EditMessageText(&msg, "获取文件失败")
		return
	}
	defer func() {
		err = os.Remove(filePath)
		if err != nil {
			logger.Info.Println("[util][upload.UploadCommand]删除文件失败:", err)
		}
	}()

	objKey, err := servcie.CosUpload(filePath, file.SavePath, file.FileName)
	if err != nil {
		logger.Info.Println("[util][upload.UploadCommand]上传到cos失败:", err)
		util.EditMessageText(&msg, "上传至cos失败")
		return
	}

	//键盘
	r := util.RandInt()
	sign := util.Sha256Hex([]byte(fmt.Sprintf("del#%s#%d", objKey, r)))[:6]
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("查看", servcie.GetFileCdnUrl(objKey)),
			tgbotapi.NewInlineKeyboardButtonData("删除", "del"),
		),
	)
	cdnUrl := ""
	qrcodeNum := 1
	if config.Get().Cos.CdnUrlDomain != "" {
		cdnUrl = "\n下载地址(CDN):\n" + servcie.GetFileCdnUrl(objKey)
		qrcodeNum = 2
	}

	//发送二维码
	qrcodes := make([]interface{}, qrcodeNum)
	qr1, _ := qrcode.Encode(servcie.GetFileUrl(objKey), qrcode.Medium, 256)
	qrcodes[0] = tgbotapi.NewInputMediaPhoto(tgbotapi.FileBytes{Name: "normal.png", Bytes: qr1})
	if config.Get().Cos.CdnUrlDomain != "" {
		qr2, _ := qrcode.Encode(servcie.GetFileUrl(objKey), qrcode.Medium, 256)
		qrcodes[1] = tgbotapi.NewInputMediaPhoto(tgbotapi.FileBytes{Name: "cdn.png", Bytes: qr2})
	}
	picMsgs, err := util.SendPics(msg.Chat.ID, qrcodes, msg.MessageID)
	if err != nil {
		logger.Error.Println("[util][upload.UploadCommand]发送二维码信息失败：", err)
		return
	}
	picMsgIDText := ""
	for _, m := range picMsgs {
		if picMsgIDText != "" {
			picMsgIDText += ","
		}
		picMsgIDText += fmt.Sprintf("%d", m.MessageID)
	}

	//编辑原消息
	err = util.EditMessageWithMarkUP(&msg, fmt.Sprintf(
		"【上传成功】\nObjectKey:【%s】 \nKey:【%s】\nQrcode:【%s】\n下载地址：\n%s%s",
		objKey,
		fmt.Sprintf("%d#%s", r, sign),
		picMsgIDText,
		servcie.GetFileUrl(objKey),
		cdnUrl,
	),
		&keyboard)
	if err != nil {
		util.EditMessageText(&msg, "异常:"+err.Error())
		return
	}
	return
}
