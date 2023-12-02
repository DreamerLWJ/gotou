package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_testHooker = ""
	_testSecret = ""
)

func TestFeishuHooker_SendTextMsg(t *testing.T) {
	hooker := NewFeishuHooker(_testHooker, Secret(_testSecret), AtAll())
	err := hooker.SendTextMsg("测试普通文本告警")
	assert.Nil(t, err)
}

func TestFeishuHooker_SendRichTextMsg(t *testing.T) {
	hooker := NewFeishuHooker(_testHooker,
		Secret(_testSecret), AtAll())
	items := make([][]RichContentItem, 0, 1)
	items = append(items, []RichContentItem{
		{
			Tag:  "text",
			Text: "服务进程异常",
		},
		{
			Tag:  "a",
			Text: "请查看",
			Href: "www.baidu.com",
		},
		{
			Tag:  "at",
			Text: "all",
		},
	})
	err := hooker.SendRichTextMsg(RichContent{
		Title:   "测试富文本告警",
		Content: items,
	})
	assert.Nil(t, err)
}
