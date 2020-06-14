package main

import (
	"os"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func initRouter() {

	port := os.Getenv("PORT")
	router := gin.Default()

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// API
	router.GET("", defaultHandler)
	router.GET("/time", timeHandler(bot))
	router.GET("/final", finalHandler(bot))
	router.POST("/webhook", lineHandler(bot))

	rg := router.Group("/report")
	{
		rg.GET("/1000", reportHandler(bot, "1000"))
		rg.GET("/1500", reportHandler(bot, "1500"))
		rg.GET("/1900", reportHandler(bot, "1900"))
	}
		
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}

	return
}

func defaultHandler(c *gin.Context) {
	i, _ := __redis.Incr("counter")
	c.String(200, "Visit number: %d", i)
}

func timeHandler(bot *linebot.Client) gin.HandlerFunc {
	
	return func(c *gin.Context) {

		taipei, _ := time.LoadLocation("Asia/Taipei")
		pushMessage := linebot.NewTextMessage(fmt.Sprintf("current time is: %v", time.Now().In(taipei)))
		if _, err := bot.PushMessage(MyID, pushMessage).Do(); err != nil {
			log.Print(err)
		}
	}
}