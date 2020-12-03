package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"alertmanager-wechatrobot-webhook/model"
	"alertmanager-wechatrobot-webhook/notifier"
)

var (
	h        bool
	RobotKey string
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.StringVar(&RobotKey, "RobotKey", "", "global wechatrobot webhook, you can overwrite by alert rule with annotations wechatRobot")
}

func main() {

	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	router := gin.Default()
	router.POST("/webhook", func(c *gin.Context) {
		var notification model.Notification

		err := c.BindJSON(&notification)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		RobotKey := c.DefaultQuery("key", RobotKey)

		err = notifier.Send(notification, RobotKey)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}

		c.JSON(http.StatusOK, gin.H{"message": "send to wechatbot successful!"})

	})
	err := router.Run(":8999")
	if err != nil{
		panic(err)
	}
}