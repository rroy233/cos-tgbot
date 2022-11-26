package util

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/logger"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var bot *tgbotapi.BotAPI

func InitUtil(b *tgbotapi.BotAPI) {
	bot = b

	_, err := os.Stat("./storage/")
	if err != nil {
		if os.IsNotExist(err) == true {
			//创建目录storage
			err = os.Mkdir("./storage", 0755)
			if err != nil {
				logger.Error.Fatalln(err)
			}
		}
	}
}

func SendPlainText(update *tgbotapi.Update, text string) {
	if update.Message == nil {
		return
	}
	var msg tgbotapi.MessageConfig
	var err error
	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(msg)
	} else if update.CallbackQuery != nil || update.CallbackQuery.Message != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		_, err = bot.Send(msg)
	}
	if err != nil {
		logger.Info.Println("[util][SendPlainText]" + err.Error())
	}
}

func DownloadFile(fileUrl string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fileUrl, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//get file name
	oFileName := "file"
	urls := strings.Split(fileUrl, "/")
	if len(urls) == 0 {
		return "", errors.New("url无效")
	}
	if strings.Contains(urls[len(urls)-1], ".") != false {
		oFileName = urls[len(urls)-1]
	}

	fileName := fmt.Sprintf("./storage/upload_%d_%s", time.Now().UnixMicro(), oFileName)
	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func EditMessageText(msg *tgbotapi.Message, newText string) {
	edit := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, newText)
	if _, err := bot.Send(edit); err != nil {
		logger.Info.Println("[util][EditMessageText]", err)
	}
}

func EditMessageWithMarkUP(msg *tgbotapi.Message, newText string, keyboard *tgbotapi.InlineKeyboardMarkup) error {
	edit := tgbotapi.NewEditMessageTextAndMarkup(msg.Chat.ID, msg.MessageID, newText, *keyboard)
	_, err := bot.Send(edit)
	if err != nil {
		logger.Info.Println("[util][EditMessageText]", err)
	}
	return err
}

func CallBack(callbackQueryID string, text string) error {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	//不能用bot.Send(callback)方法，有bug
	resp, err := bot.Request(callback)
	if err != nil {
		return err
	}
	if string(resp.Result) != "true" {
		return errors.New("请求不ok")
	}
	return err
}
func CallBackWithAlert(callbackQueryID string, text string) error {
	callback := tgbotapi.NewCallbackWithAlert(callbackQueryID, text)
	//不能用bot.Send(callback)方法，有bug
	resp, err := bot.Request(callback)
	if err != nil {
		return err
	}
	if string(resp.Result) != "true" {
		return errors.New("请求不ok")
	}

	return err
}

// RandInt 生成6位随机数
func RandInt() int {
	rand.Seed(time.Now().UnixMilli())
	return 100000 + rand.Intn(899999)
}

func JsonEncode(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
func MatchSingle(re *regexp.Regexp, content string) (string, error) {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return "", errors.New("errorNoMatched")
	}
	return matched[0][1], nil
}
