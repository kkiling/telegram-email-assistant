package email

import "github.com/kiling91/telegram-email-assistant/pkg/types"

type ImapServer interface {
	ReadUnseenEmails(user *types.EmailUser) error
}
