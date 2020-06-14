package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	DefaultMessage = "在家休息 無感冒發燒"
	ReturnMessage = "在家休假 預計1730至中正紀念堂"
	TaoyuanMessage = "在家休假 預計17到公園集合"
	TitleMessage = "今日看診人員:共0員\n發燒人員:共0員\n應到：15員\n實到：15員\n"
)

func lineHandler(bot *linebot.Client) gin.HandlerFunc {

	return func(c *gin.Context) {

		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.String(400, "%v", err)
			} else {
				c.String(500, "%v", err)
			}
			return
		}
		
		for _, event := range events {
			switch event.Type {
				
			case linebot.EventTypeFollow:
				followRoutine(bot, event)

			case linebot.EventTypeJoin:
				joinRoutine(bot, event)

			case linebot.EventTypePostback:
				postbackRoutine(bot, event)

			case linebot.EventTypeMessage:
				messageRoutine(bot, event)
			}
		}
	}
}

func reportHandler(bot *linebot.Client, clock string) gin.HandlerFunc {
	
	return func(c *gin.Context) {

		date := getDateTime(clock)
		reportMessage := linebot.NewTextMessage(date + TitleMessage + getCombinedMessage())
		if _, err := bot.PushMessage(CompanyID, reportMessage).Do(); err != nil {
			log.Print(err)
		}
	}
}

func finalHandler(bot *linebot.Client) gin.HandlerFunc {

	return func(c *gin.Context) {

		studentIDs, err := __redis.SMembers("students")
		if err != nil {
			log.Print(err)
			return
		}
	
		sort.Strings(studentIDs)
		
		var finalMsg string
		for _, studentID := range studentIDs {
	
			userID, _, err := __redis.Get(studentID)
			if err != nil{
				log.Print(err)
				continue
			}
	
			user := GetUser(userID)
			if user == nil {
				continue
			}
	
			info := fmt.Sprintf("%s %s %s\n", user.StudentID, user.Name, user.PhoneNumber)
			/*
			var msg string
			if studentID == "33090" {
				msg = TaoyuanMessage
			} else {
				msg = ReturnMessage
			}

			if user.Message != "" {
				msg = user.Message
			}
			*/
			msg := user.Message
			user.Message = ""
			if err := __redis.Set(userID, user.ToJSON(), 0); err != nil {
				log.Print(err)
				continue
			}
	
			finalMsg = finalMsg + info + msg + "\n\n"
		}

		date := getDateTime("1500")
		reportMessage := linebot.NewTextMessage(date + TitleMessage + finalMsg)
		if _, err := bot.PushMessage(CompanyID, reportMessage).Do(); err != nil {
			log.Print(err)
		}
	}
}

func getCombinedMessage() (finalMsg string) {

	studentIDs, err := __redis.SMembers("students")
	if err != nil {
		log.Print(err)
		return
	}

	sort.Strings(studentIDs)
	
	for _, studentID := range studentIDs {

		userID, _, err := __redis.Get(studentID)
		if err != nil{
			log.Print(err)
			continue
		}

		user := GetUser(userID)
		if user == nil {
			continue
		}

		info := fmt.Sprintf("%s %s %s\n", user.StudentID, user.Name, user.PhoneNumber)
		/*
		msg := DefaultMessage
		if user.Message != "" {
			msg = user.Message
		}
		*/
		msg := user.Message
		user.Message = ""
		if err := __redis.Set(userID, user.ToJSON(), 0); err != nil {
			log.Print(err)
			continue
		}

		finalMsg = finalMsg + info + msg + "\n\n"
	}

	return
}

func getDateTime(clock string) (dateMsg string) {

	dateMsg = "109/" + time.Now().Format("01/02")

	switch clock {
	case "1000": dateMsg = dateMsg + "休假10點回報\n"
	case "1500": dateMsg = dateMsg + "休假15點回報\n"
	case "1900": dateMsg = dateMsg + "休假19點回報\n"
	}

	return
}