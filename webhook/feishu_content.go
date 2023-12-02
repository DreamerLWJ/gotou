package webhook

import "fmt"

const (
	_hookTextRefUserContentFormat = "<at user_id=\\\"%s\\\">%s</at> %s"
)

func getHookTextRefUser(userID, refUserText, content string) string {
	return fmt.Sprintf(_hookTextRefUserContentFormat, userID, refUserText, content)
}

type feishuEncryptFields struct {
	Timestamp string `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
}

type feishuTextHookReq struct {
	feishuEncryptFields
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type RichContent struct {
	Title   string              `json:"title"`
	Content [][]RichContentItem `json:"content"`
}

type RichContentItem struct {
	Tag    string `json:"tag"`
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type feishuRichTextHookReq struct {
	feishuEncryptFields
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn RichContent `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

type feishuHookResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
