package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type dingdingTextMsgResp struct {
	Errcode int
	Errmsg  string
}
type dingdingTextMsg struct {
	MsgType string         `json:"msgtype"`
	Text    dingtextParams `json:"text"`
	At      dingdingAt     `json:"at"`
}

type dingtextParams struct {
	Content string `json:"content"`
}

// At at struct
type dingdingAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func SendDingRobot(url string, content string, isAtAll bool) error {
	msg := dingdingTextMsg{MsgType: "text", Text: dingtextParams{Content: content}, At: dingdingAt{IsAtAll: isAtAll}}
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(m))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dr dingdingTextMsgResp
	err = json.Unmarshal(data, &dr)
	if err != nil {
		return err
	}
	if dr.Errcode != 0 {
		return fmt.Errorf("dingrobot send failed: %v", dr.Errmsg)
	}
	return nil
}

type dingdingMarkdownReq struct {
	MsgType  string      `json:"msgtype"`
	Markdown markdownMsg `json:"markdown"`
	At       dingdingAt  `json:"at"`
}

type markdownMsg struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// 支持Markdown格式的钉钉告警机器人
// 注意 在text内容里要有@人的手机号,且与atMobiles中的手机号一一对应才有用
func SendDingRobotMD(url string, title, text string, isAtAll bool, atMobiles ...string) error {
	msg := dingdingMarkdownReq{
		MsgType:  "markdown",
		Markdown: markdownMsg{Title: title, Text: text},
		At:       dingdingAt{IsAtAll: isAtAll, AtMobiles: atMobiles},
	}
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(m))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dr dingdingTextMsgResp
	err = json.Unmarshal(data, &dr)
	if err != nil {
		return err
	}
	if dr.Errcode != 0 {
		return fmt.Errorf("dingrobot send failed: %v", dr.Errmsg)
	}
	return nil
}

type SendDingReportReq struct {
	Url     string   `json:"url"`
	Content string   `json:"content"`
	All     bool     `json:"all"`
	Phones  []string `json:"phones"`
}

// TextMessage 发送钉钉消息
type TextMessage struct {
	MsgType string     `json:"msgtype"`
	Text    TextParams `json:"text"`
	At      struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

type TextParams struct {
	Content string `json:"content"`
}

type DingResponse struct {
	Errcode int
	Errmsg  string
}

// 钉钉告警
func DingDingDisturbSpecifiedStaff(ctx context.Context, req SendDingReportReq) error {
	logHead := "DingDingDisturbSpecifiedStaff|"
	msgBody := TextMessage{MsgType: "text", Text: TextParams{Content: req.Content}, At: struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	}(struct {
		AtMobiles []string
		IsAtAll   bool
	}{AtMobiles: req.Phones, IsAtAll: req.All})}
	m, err := json.Marshal(msgBody)
	if err != nil {
		return err
	}
	resp, err := http.Post(req.Url, "application/json", bytes.NewReader(m))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	var dr DingResponse
	err = json.Unmarshal(data, &dr)
	if err != nil {
		return errors.WithStack(err)
	}
	if dr.Errcode != 0 {
		return fmt.Errorf(logHead+"dingrobot send failed: %v", dr.Errmsg)
	}
	return nil
}
