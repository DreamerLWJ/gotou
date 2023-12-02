package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type FeishuHooker struct {
	url     string
	atAll   bool
	secret  string
	keyword string
}

func (f *FeishuHooker) setSecret(s string) {
	f.secret = s
}

func (f *FeishuHooker) setKeyword(k string) {
	f.keyword = k
}

func NewFeishuHooker(url string, opts ...Option) *FeishuHooker {
	hooker := &FeishuHooker{url: url}
	for _, opt := range opts {
		opt(hooker)
	}
	return hooker
}

func (f *FeishuHooker) processEncrypt() (feishuEncryptFields, error) {
	if f.secret == "" {
		return feishuEncryptFields{}, nil
	}

	sign, err := f.GenSign(f.secret, strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return feishuEncryptFields{}, errors.Errorf("FeishuHooker|gen sign err:%s", err)
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	return feishuEncryptFields{
		Timestamp: ts,
		Sign:      sign,
	}, nil
}

func (f *FeishuHooker) SendTextMsg(content string) error {
	encryptFlds, err := f.processEncrypt()
	if err != nil {
		return err
	}
	req := feishuTextHookReq{
		feishuEncryptFields: encryptFlds,
		MsgType:             "text",
		Content: struct {
			Text string `json:"text"`
		}{content},
	}
	return f.doSend(req)
}

func (f *FeishuHooker) SendRichTextMsg(content RichContent) error {
	encryptFlds, err := f.processEncrypt()
	if err != nil {
		return err
	}
	req := feishuRichTextHookReq{
		feishuEncryptFields: encryptFlds,
		MsgType:             "post",
		Content: struct {
			Post struct {
				ZhCn RichContent `json:"zh_cn"`
			} `json:"post"`
		}{
			Post: struct {
				ZhCn RichContent `json:"zh_cn"`
			}{ZhCn: content},
		},
	}
	err = f.doSend(req)
	if err != nil {
		return err
	}
	return nil
}

func (f *FeishuHooker) refAllUser() {
	f.atAll = true
}

func (f *FeishuHooker) GenSign(secret string, timestamp string) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := timestamp + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func (f *FeishuHooker) doSend(req any) error {

	data, _ := json.Marshal(req)
	response, err := http.Post(f.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return errors.Errorf("FeishuHooker|post err:%s", err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			// TODO
		}
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	resp := feishuHookResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return errors.Errorf("FeishuHooker|read resp err:%s", err)
	}
	if resp.Code != 0 || resp.Msg != "success" {
		return errors.Errorf("FeishuHooker|resp err (code:%d,msg:%s)", resp.Code, resp.Msg)
	}
	return nil
}
