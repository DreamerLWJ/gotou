package webhook

type SupportHookType int8

const (
	Feishu   SupportHookType = 1
	DingDing SupportHookType = 2
)

type Option func(p Hooker)

func AtAll() Option {
	return func(p Hooker) {
		p.refAllUser()
	}
}

func Secret(sc string) Option {
	return func(p Hooker) {
		p.setSecret(sc)
	}
}

func Keyword(k string) Option {
	return func(p Hooker) {
		p.setKeyword(k)
	}
}

type Hooker interface {
	SendTextMsg(content string) error

	SendRichTextMsg(content RichContent) error

	refAllUser()

	setSecret(s string)

	setKeyword(k string)
}
