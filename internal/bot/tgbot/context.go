package tgbot

import (
	"github.com/kiling91/telegram-email-assistant/internal/bot"
	tg "gopkg.in/telebot.v3"
)

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

// ***
type btnContext struct {
	ctx tg.Context
}

func newBtnContext(ctx tg.Context) bot.BtnContext {
	return &btnContext{
		ctx: ctx,
	}
}

func (c *btnContext) UserId() int64 {
	return c.ctx.Sender().ID
}

func (c *btnContext) Unique() string {
	return c.ctx.Callback().ID
}

func (c *btnContext) Data() string {
	return c.ctx.Callback().Data
}
