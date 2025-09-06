package transformer

import (
	"alertmanager-wechatrobot-webhook/model"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"time"
)

// TransformToMarkdown transform alertmanager notification to wechat markdow message
func TransformToMarkdown(notification model.Notification) (markdown *model.WeChatMarkdown, robotURL string, err error) {

	annotations := notification.CommonAnnotations
	robotURL = annotations["wechatRobot"]

	// 如果告警数量较多，尝试拆分发送
	if len(notification.Alerts) > 1 {
		return transformWithSplit(notification, robotURL)
	}

	var buffer bytes.Buffer

	content, err := toContent(notification)
	buffer.WriteString(content)

	// 限制内容长度，企业微信markdown内容最大4096字符
	finalContent := buffer.String()
	if len(finalContent) > 4000 { // 留一些余量
		finalContent = finalContent[:4000] + "\n\n...内容过长已截断"
	}

	markdown = &model.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &model.Markdown{
			Content: finalContent,
		},
	}

	return
}

// transformWithSplit 当告警内容过长时，拆分成多个消息
func transformWithSplit(notification model.Notification, robotURL string) (markdown *model.WeChatMarkdown, url string, err error) {
	// 先尝试生成完整内容
	content, err := toContent(notification)
	if err != nil {
		return nil, robotURL, err
	}

	// 如果内容长度在限制内，直接返回
	if len(content) <= 4000 {
		markdown = &model.WeChatMarkdown{
			MsgType: "markdown",
			Markdown: &model.Markdown{
				Content: content,
			},
		}
		return markdown, robotURL, nil
	}

	// 内容过长，生成第一个告警的消息
	firstAlert := model.Notification{
		Version:           notification.Version,
		GroupKey:          notification.GroupKey,
		Status:            notification.Status,
		Receiver:          notification.Receiver,
		GroupLabels:       notification.GroupLabels,
		CommonLabels:      notification.CommonLabels,
		CommonAnnotations: notification.CommonAnnotations,
		ExternalURL:       notification.ExternalURL,
		Alerts:            []model.Alert{notification.Alerts[0]},
	}

	firstContent, err := toContent(firstAlert)
	if err != nil {
		return nil, robotURL, err
	}

	// 添加提示信息
	if len(notification.Alerts) > 1 {
		firstContent += "\n\n📢 **注意**: 本次共有 " + fmt.Sprintf("%d", len(notification.Alerts)) + " 个告警，将分批发送"
	}

	markdown = &model.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &model.Markdown{
			Content: firstContent,
		},
	}

	return markdown, robotURL, nil
}

func toContent(notification model.Notification) (content string, err error) {

	templateString := templateString()
	out := bytes.NewBuffer([]byte{})

	funcMap := template.FuncMap{"fdate": formDate}
	t, err := template.New("test").Funcs(funcMap).Parse(templateString)
	if err != nil {
		return
	}
	err = t.Execute(out, notification)
	return out.String(), nil
}

func formDate(t time.Time) string {
	layout := "2006-01-02 15:04:05"
	// loc, err := time.LoadLocation("Local")
	// if err != nil {
	// 	log.Panicf("LoadLocation err %v", err)
	// 	return t.Format(layout)
	// }
	return t.Local().Format(layout)
}

var defaultTemplateString = `# {{ if eq .Status "resolved"}}<font color="info">恢复</font>{{ else if eq .Status "firing"}}<font color="warning">触发</font>{{end}}  {{.CommonLabels.alertname}}  
{{if .CommonLabels.severity          }}## 级别: <font color="warning">{{.CommonLabels.severity}}</font>  {{end}} 
{{ if .CommonAnnotations.description }}## 描述: {{.CommonAnnotations.description }} {{end}}
{{ if .CommonAnnotations.summary     }}## 汇总: {{.CommonAnnotations.summary}}      {{end}}
{{ range .Alerts}}
##### 标签: 
{{ range $key, $value := .Labels}}
1. {{$key}}: {{$value}}  
{{end}}

{{ if not .StartsAt.IsZero }}触发时间 {{.StartsAt | fdate}} {{end}} 
{{ if not .EndsAt.IsZero }}恢复时间 {{.EndsAt | fdate}}   {{end}}
{{end}}
`

func templateString() string {
	filePath := os.Getenv("template_path")
	if filePath == "" {
		return defaultTemplateString
	}
	file, err := os.Open(filePath)
	if err != nil {
		return defaultTemplateString
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return defaultTemplateString
	}
	return string(b)
}
