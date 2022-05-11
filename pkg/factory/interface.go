package factory

import (
	"github.com/kiling91/telegram-email-assistant/pkg/email"
	"github.com/kiling91/telegram-email-assistant/pkg/storage"
)

type Factory interface {
	GetStorage() storage.Storage
	ImapServer() email.ImapServer
}
