package tgbot

import (
	"log"

	tg "gopkg.in/telebot.v3"
)

func (t *telegram) drawInlineBtn(inlineBtns []*inlineBtn) interface{} {
	menu := &tg.ReplyMarkup{}

	rows := make([]tg.Row, len(inlineBtns))
	for i, btn := range inlineBtns {
		b := menu.Data(btn.text, btn.unique)

		t.bot.Handle(&b, func(c tg.Context) error {
			log.Println(c.Callback().Unique)
			return c.Respond()
		})

		rows[i] = menu.Row(b)
	}

	menu.Inline(rows...)
	return menu
}

func (t *telegram) extractOptions(how []interface{}) []interface{} {

	var opts = make([]interface{}, len(how))
	for i, prop := range how {
		switch opt := prop.(type) {
		case []*inlineBtn:
			opts[i] = t.drawInlineBtn(opt)
		default:
			panic("telebot: unsupported send-option")
		}
	}

	return opts
}
