package factory

import (
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/storage"
)

type Factory interface {
	GetStorage() storage.Storage
	ImapServer() email.ImapServer
}
