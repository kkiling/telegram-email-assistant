package tgbot

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/kiling91/telegram-email-assistant/pkg/bot"
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

func NewTbBot(token string) (bot.Bot, error) {
	pref := tg.Settings{
		Token:  token,
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

func (t *telegram) SendPhoto(userId int64, photo *bot.Photo, opts ...interface{}) (e *bot.Editable, err error) {
	var msg []tg.Message

	p := &tg.Photo{
		File:    tg.FromDisk(photo.Filename),
		Caption: photo.Caption,
	}

	if len(opts) > 0 {
		sendOptions := t.extractOptions(opts)
		msg, err = t.bot.SendAlbum(newUser(userId), tg.Album{p}, sendOptions)
	} else {
		msg, err = t.bot.SendAlbum(newUser(userId), tg.Album{p})
	}

	if err != nil {
		return nil, err
	}

	return &bot.Editable{
		MessageID: msg[0].ID,
		ChatID:    msg[0].Chat.ID,
	}, err
}

func (t *telegram) Send(userId int64, text string, opts ...interface{}) (e *bot.Editable, err error) {
	var msg *tg.Message

	if len(opts) > 0 {
		sendOptions := t.extractOptions(opts)
		msg, err = t.bot.Send(newUser(userId), text, sendOptions, tg.ModeHTML)
	} else {
		msg, err = t.bot.Send(newUser(userId), text, tg.ModeHTML)
	}

	if err != nil {
		return nil, err
	}

	return &bot.Editable{
		MessageID: msg.ID,
		ChatID:    msg.Chat.ID,
	}, err
}

func (t *telegram) SendDocument(userId int64, filename string) error {
	doc := &tg.Document{
		File:     tg.FromDisk(filename),
		FileName: filepath.Base(filename),
	}
	_, err := t.bot.Send(newUser(userId), doc)
	return err
}

func (t *telegram) Edit(edit *bot.Editable, text string, opts ...interface{}) (e *bot.Editable, err error) {
	var msg *tg.Message

	if len(opts) > 0 {
		sendOptions := t.extractOptions(opts)
		msg, err = t.bot.Edit(edit, text, sendOptions, tg.ModeHTML)
	} else {
		msg, err = t.bot.Edit(edit, text, tg.ModeHTML)
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
