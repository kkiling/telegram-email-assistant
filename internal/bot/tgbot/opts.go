package tgbot

import (
	"github.com/kiling91/telegram-email-assistant/internal/bot"
	tg "gopkg.in/telebot.v3"
)

func (t *telegram) drawInlineBtn(inline *bot.Inline) *tg.ReplyMarkup {
	menu := &tg.ReplyMarkup{}

	btns := inline.GetBtns()

	var rows []tg.Row
	var row []tg.Btn

	for index, btn := range btns {
		if index%inline.ItemsPerRow == 0 {
			if len(row) > 0 {
				rows = append(rows, menu.Row(row...))
			}
			row = []tg.Btn{}
		}
		b := menu.Data(btn.Text, btn.Unique, btn.Data)

		t.bot.Handle(&b, func(c tg.Context) error {
			return inline.Handler(newBtnContext(c))
		})

		row = append(row, b)
	}

	if len(row) > 0 {
		rows = append(rows, menu.Row(row...))
	}
	menu.Inline(rows...)
	return menu
}

func (t *telegram) extractOptions(how []interface{}) *tg.SendOptions {
	opts := &tg.SendOptions{}

	for _, prop := range how {
		switch opt := prop.(type) {
		case *bot.Inline:
			opts.ReplyMarkup = t.drawInlineBtn(opt)
		default:
			panic("telebot: unsupported send-option")
		}
	}

	return opts
}
