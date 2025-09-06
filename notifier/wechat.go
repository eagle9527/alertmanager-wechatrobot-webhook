package notifier

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"alertmanager-wechatrobot-webhook/model"
	"alertmanager-wechatrobot-webhook/transformer"
)

// Send send markdown message to dingtalk
func Send(notification model.Notification, defaultRobot string) (err error) {

	markdown, robotURL, err := transformer.TransformToMarkdown(notification)

	if err != nil {
		return
	}

	SendMarkDown(markdown, robotURL, defaultRobot)

	// 如果有多个告警且内容可能被拆分，发送剩余告警
	if len(notification.Alerts) > 1 {
		err = sendRemainingAlerts(notification, robotURL, defaultRobot)
	}

	return
}

// sendRemainingAlerts 发送剩余的告警
func sendRemainingAlerts(notification model.Notification, robotURL string, defaultRobot string) error {
	// 从第二个告警开始发送
	for i := 1; i < len(notification.Alerts); i++ {
		// 创建单个告警的通知
		singleAlert := model.Notification{
			Version:           notification.Version,
			GroupKey:          notification.GroupKey,
			Status:            notification.Status,
			Receiver:          notification.Receiver,
			GroupLabels:       notification.GroupLabels,
			CommonLabels:      notification.CommonLabels,
			CommonAnnotations: notification.CommonAnnotations,
			ExternalURL:       notification.ExternalURL,
			Alerts:            []model.Alert{notification.Alerts[i]},
		}

		markdown, _, err := transformer.TransformToMarkdown(singleAlert)
		if err != nil {
			log.Printf("Error transforming alert %d: %v", i+1, err)
			continue
		}

		SendMarkDown(markdown, robotURL, defaultRobot)
	}
	return nil
}

func SendMarkDown(markdown *model.WeChatMarkdown, robotURL string, robot string) {
	data, err := json.Marshal(markdown)

	println(data)
	if err != nil {
		return
	}

	var wechatRobotURL string
	if robotURL != "" {
		wechatRobotURL = robotURL
	} else {
		wechatRobotURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + robot
	}

	req, err := http.NewRequest(
		"POST",
		wechatRobotURL,
		bytes.NewBuffer(data))

	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	tr := &http.Transport{
		//TLSClientConfig:        &tls.Config{InsecureSkipVerify:true},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	return
}
