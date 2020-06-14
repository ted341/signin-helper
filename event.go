package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func followRoutine(bot *linebot.Client, event *linebot.Event) {

	userID := event.Source.UserID
	profile, err := bot.GetProfile(userID).Do()
	if err != nil {
		log.Print(err)
		return
	}

	displayName := profile.DisplayName
	idx := regexp.MustCompile("^[0-9]*").FindStringIndex(displayName)
	id  := displayName[idx[0]:idx[1]]
	name := strings.TrimSpace(displayName[idx[1]:len(displayName)])
	
	if id == "" {
		msg := linebot.NewTextMessage("請在你的名字前面加上學號好讓我認識你唷~(ex.33076張翔中)")
		if _, err := bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
			log.Print(err)
		}
		return
	}
	
	user := &User{StudentID: id}
	if p, ok := MySquad[id]; !ok {
		msg := linebot.NewTextMessage("抱歉，你不是六班的訓員，所以我無法幫助你喔~")
		if _, err := bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
			log.Print(err)
		}
		return

	} else {
		user.Profile = p
	}

	if err := __redis.Set(userID, user.ToJSON(), 0); err != nil {
		log.Print(err)
	}

	if err := __redis.Set(id, userID, 0); err != nil {
		log.Print(err)
	}

	if err := __redis.SAdd("students", id); err != nil {
		log.Print(err)
	}

	replyMessage := linebot.NewTextMessage(fmt.Sprintf("hi, %s", name))
	if _, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
		log.Print(err)
	}
}

func joinRoutine(bot *linebot.Client, event *linebot.Event) {

	pushMessage := linebot.NewTextMessage(fmt.Sprintf("Joined %s", event.Source.GroupID))
	if _, err := bot.PushMessage(MyID, pushMessage).Do(); err != nil {
		log.Print(err)
	}
}

func postbackRoutine(bot *linebot.Client, event *linebot.Event) {

	userID := event.Source.UserID

	switch event.Postback.Data {
	case "report":
		user := GetUser(userID)
		if user == nil {
			return
		}

		user.Unblocked = true
		if err := __redis.Set(userID, user.ToJSON(), 0); err != nil {
			log.Print(err)
			return
		}

		replyMessage := linebot.NewTextMessage("請開始回報~")
		if _, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
			log.Print(err)
		}

	case "inquire":
		msg := getCombinedMessage()
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
			log.Print(err)
		}
	}
}

func messageRoutine(bot *linebot.Client, event *linebot.Event) {

	userID := event.Source.UserID
	user := GetUser(userID)
	if user == nil || user.Unblocked != true {
		return
	}

	if text, ok := event.Message.(*linebot.TextMessage); ok {
		user.Message = text.Text
		user.Unblocked = false
		if err := __redis.Set(userID, user.ToJSON(), 0); err != nil {
			log.Print(err)
			return
		}
		replyMessage := linebot.NewTextMessage("收到~")
		if _, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
			log.Print(err)
		}

	} else {
		replyMessage := linebot.NewTextMessage("請回報文字訊息~")
		if _, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
			log.Print(err)
		}
	}

	return
}