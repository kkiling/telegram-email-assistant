package email

import "time"

type ImapUser struct {
	ImapServer string
	Login      string
	Password   string
}

type Message struct {
	// The message unique identifier. It must be greater than or equal to 1.
	Uid uint32
	// The message date.
	Date time.Time
	// The message subject.
	Subject string
	// From header addresses.
	From string
	// The message senders.
	To string
}

type InlineFile struct {
	FileName string
	FilePath string
}

type AttachmentFile struct {
	FileName string
	FilePath string
}

type MessageWithBody struct {
	Message
	TextPlain       string
	TextHtml        string
	InlineFiles     []*InlineFile
	AttachmentFiles []*AttachmentFile
}
