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

type context struct {
	ctx tg.Context
}

func newContext(ctx tg.Context) bot.Context {
	return &context{
		ctx: ctx,
	}
}

func (c *context) UserId() int64 {
	return c.ctx.Sender().ID
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

func (t *telegram) Send(userId int64, text string, opts ...interface{}) (msgId int, err error) {
	msg, err := t.bot.Send(newUser(userId), text, t.extractOptions(opts))
	if err != nil {
		return 0, err
	}
	return msg.ID, err
}
