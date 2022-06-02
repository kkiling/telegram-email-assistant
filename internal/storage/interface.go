package storage

import "github.com/kiling91/telegram-email-assistant/internal/email"

type Storage interface {
	// Отправленные непрочитанные сообщения пользователю бота
	SaveMsgSentToBotUser(email string, seqNum int64, botUserUid int64) error
	// Проверяем было ли сообщение уже отправленно пользователю бота
	MsgWasSentToBotUser(email string, seqNum int64, botUserUid int64) (bool, error)
	//
	SaveMsgInfo(email string, msg *email.MessageEnvelope) error
	GetMsgFromAddress(email string, seqNum int64) (string, error)
	ShutDown() error
}
