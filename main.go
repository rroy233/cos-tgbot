package main

import (
	"context"
	"encoding/json"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"tgUploader/command"
	"tgUploader/config"
	"tgUploader/servcie"
	"tgUploader/util"
)

const WorkerNum = 1

var bot *tgbotapi.BotAPI

func main() {
	var debug *bool
	debug = flag.Bool("debug", false, "debug mode")
	flag.Parse()

	//加载配置文件
	config.LoadConfig()

	var err error
	bot, err = tgbotapi.NewBotAPI(config.Get().BotToken)
	if err != nil {
		log.Panic(err)
	}

	command.InitCommand(bot)
	util.InitUtil(bot)
	servcie.InitService()

	bot.Debug = *debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	stopCtx, cancel := context.WithCancel(context.Background())
	cancelCh := make(chan int, WorkerNum)
	for i := 0; i < WorkerNum; i++ {
		go worker(stopCtx, updates, cancelCh)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh

	log.Println("正在关闭服务。。。")
	cancel()
	waitForDone(cancelCh)
	log.Println("已关闭服务")
}

func worker(stopCtx context.Context, uc tgbotapi.UpdatesChannel, cancelCh chan int) {
	for {
		select {
		case update := <-uc:
			handle(update)
		case <-stopCtx.Done():
			cancelCh <- 1
			return
		}
	}
}

func handle(update tgbotapi.Update) {
	//debug
	jsonlog, _ := json.Marshal(update)
	log.Println("[update]" + string(jsonlog))

	if update.Message == nil && update.CallbackQuery == nil {
		return
	}
	command.Handle(update)
}

func waitForDone(cancelCh chan int) {
	num := 0
	for {
		if num == WorkerNum {
			break
		}
		<-cancelCh
		num++
	}
}
