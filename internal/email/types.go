package email

import "time"

type ImapUser struct {
	ImapServer string
	Login      string
	Password   string
}

type InlineFile struct {
	AttachmentId string
	FileName     string
	FilePath     string
}

type AttachmentFile struct {
	FileName string
	FilePath string
}

type Message struct {
	Envelope *MessageEnvelope
	Body     *MessageBody
}

type MessageEnvelope struct {
	// The message unique identifier. It must be greater than or equal to 1.
	SeqNum int64
	// The message date.
	Date time.Time
	// The message subject.
	Subject string
	// From header addresses.
	FromName    string
	FromAddress string
	// The message senders.
	ToName    string
	ToAddress string
}

type MessageBody struct {
	TextPlain       string
	TextHtml        string
	InlineFiles     []*InlineFile
	AttachmentFiles []*AttachmentFile
}
