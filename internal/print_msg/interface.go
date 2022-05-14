package printmsg

import "github.com/kiling91/telegram-email-assistant/internal/email"

type FormattedMsg struct {
	Text       string
	Img        string
	Attachment []string
}

type PrintMsg interface {
	PrintMsgEnvelope(msg *email.MessageEnvelope) string
	PrintMsgWithBody(msg *email.Message) (*FormattedMsg, error)
}
