package bot

import "strconv"

type UserProfile struct {
	ID int64 `json:"id"`
}

func (p *UserProfile) Recipient() string {
	return strconv.FormatInt(p.ID, 10)
}

type Editable struct {
	MessageID int   `json:"message_id"`
	ChatID    int64 `json:"chat_iD"`
}

func (e *Editable) MessageSig() (messageID string, chatID int64) {
	return strconv.Itoa(e.MessageID), e.ChatID
}

type Context interface {
	UserId() int64
}
