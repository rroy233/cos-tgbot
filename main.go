package main

import (
	"context"
	"cosTgBot/command"
	"cosTgBot/config"
	"cosTgBot/servcie"
	"cosTgBot/util"
	"encoding/json"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/logger"
	"log"
	"os"
	"os/signal"
)

const WorkerNum = 1

var bot *tgbotapi.BotAPI

func main() {
	log.Println("正在初始化")
	var debug *bool
	debug = flag.Bool("debug", false, "debug mode")
	flag.Parse()

	//加载配置文件
	config.LoadConfig()

	//日志服务
	logger.New(&logger.Config{
		StdOutput:      true,
		StoreLocalFile: config.Get().Logger.Enabled,
		StoreRemote:    config.Get().Logger.Report,
		RemoteConfig: logger.RemoteConfigStruct{
			RequestUrl: config.Get().Logger.ReportUrl,
			QueryKey:   config.Get().Logger.QueryKey,
		},
	})

	var err error
	bot, err = tgbotapi.NewBotAPI(config.Get().BotToken)
	if err != nil {
		logger.FATAL.Fatalln(err)
	}

	command.InitCommand(bot)
	util.InitUtil(bot)
	servcie.InitService()

	bot.Debug = *debug

	logger.Info.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	stopCtx, cancel := context.WithCancel(context.Background())
	cancelCh := make(chan int, WorkerNum)
	for i := 0; i < WorkerNum; i++ {
		go worker(stopCtx, updates, cancelCh)
	}

	logger.Info.Println("初始化成功")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh

	logger.Info.Println("正在关闭服务。。。")
	cancel()
	waitForDone(cancelCh)
	logger.Info.Println("已关闭服务")
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
	logger.Info.Println("[update]" + string(jsonlog))

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
