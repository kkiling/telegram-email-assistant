package bot

import "strconv"

type UserProfile struct {
	ID int64 `json:"id"`
}

func (p *UserProfile) Recipient() string {
	return strconv.FormatInt(p.ID, 10)
}

type InlineBtn struct {
	Text   string
	Unique string
	Data   string
}

type Context interface {
	UserId() int64
}

type HandlerFunc func(Context) error

type Bot interface {
	Handle(command string, h HandlerFunc)
	Send(userId int64, text string, opts ...interface{}) (msgId int, err error)
}
