package tgbot

import (
	"strconv"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/bot"
	"github.com/kiling91/telegram-email-assistant/internal/config"
	tg "gopkg.in/telebot.v3"
)

type user struct {
	id int64
}

func (u *user) Recipient() string {
	return strconv.FormatInt(u.id, 10)
}

func newUser(id int64) tg.Recipient {
	return &user{id: id}
}

type telegram struct {
	bot *tg.Bot
}

func NewTbBot(cfg *config.Telegram) (bot.Bot, error) {
	pref := tg.Settings{
		Token:  cfg.BotToken,
		Poller: &tg.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tg.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &telegram{
		bot: b,
	}, nil
}

func (t *telegram) Handle(command string, h bot.HandlerFunc) {
	t.bot.Handle(command, func(c tg.Context) error {
		return h(newContext(c))
	})
}

func (t *telegram) Send(userId int64, text string, opts ...interface{}) (e *bot.Editable, err error) {
	var msg *tg.Message

	if len(opts) > 0 {
		sendOptions := t.extractOptions(opts)
		msg, err = t.bot.Send(newUser(userId), text, sendOptions)
	} else {
		msg, err = t.bot.Send(newUser(userId), text)
	}

	if err != nil {
		return nil, err
	}

	return &bot.Editable{
		MessageID: msg.ID,
		ChatID:    msg.Chat.ID,
	}, err
}

func (t *telegram) Edit(edit *bot.Editable, text string, opts ...interface{}) (e *bot.Editable, err error) {
	var msg *tg.Message

	if len(opts) > 0 {
		sendOptions := t.extractOptions(opts)
		msg, err = t.bot.Edit(edit, text, sendOptions)
	} else {
		msg, err = t.bot.Edit(edit, text)
	}

	if err != nil {
		return nil, err
	}

	return &bot.Editable{
		MessageID: msg.ID,
		ChatID:    msg.Chat.ID,
	}, err
}

func (t *telegram) Start() {
	t.bot.Start()
}

func (t *telegram) Stop() {
	t.bot.Stop()
}
