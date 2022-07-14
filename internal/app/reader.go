package app

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/kiling91/telegram-email-assistant/internal/printmsg"
	"github.com/kiling91/telegram-email-assistant/pkg/bot"
	"github.com/sirupsen/logrus"
)

const progressTimeout = 10

type Reader struct {
	fact     factory.Factory
	userIds  []int64
	imapUser *email.ImapUser
}

func NewReader(fact factory.Factory, userIds []int64, imapUser *email.ImapUser) *Reader {
	return &Reader{
		fact:     fact,
		userIds:  userIds,
		imapUser: imapUser,
	}
}

func (r *Reader) sendPrintMsg(fmsg *printmsg.FormattedMsg, userId int64) {
	b := r.fact.Bot()

	if fmsg.Img != "" {
		_, err := b.SendPhoto(userId, &bot.Photo{
			Filename: fmsg.Img,
			Caption:  fmsg.Text,
		})
		if err != nil {
			logrus.Errorf("error send photo: %v", err)
			return
		}
	} else {
		_, err := b.Send(userId, fmsg.Text)
		if err != nil {
			logrus.Errorf("error send photo: %v", err)
			return
		}
	}

	for _, attach := range fmsg.Attachment {
		err := b.SendDocument(userId, attach)
		if err != nil {
			logrus.Errorf("error send document: %v", err)
			return
		}
	}
}

func (r *Reader) startReadProgress(ctx context.Context, userId int64, seqNum int64, end <-chan bool) {
	b := r.fact.Bot()

	storage := r.fact.Storage()
	from, err := storage.GetMsgFromAddress(r.imapUser.Login, seqNum)
	if err != nil {
		logrus.Errorf("error get msg info: %v", err)
		return
	}
	edit, err := b.Send(userId, fmt.Sprintf("⌛ Reading a mail from %s", from))
	if err != nil {
		logrus.Errorf("error send msg to user %d", userId)
		return
	}
	go func() {
		counter := 0
		for {
			timer := time.NewTimer(progressTimeout * time.Second)
			select {
			case <-ctx.Done():
				return
			case <-end:
				b.Delete(edit)
				return
			case <-timer.C:
				counter += 1
				if counter%2 == 0 {
					b.Edit(edit, fmt.Sprintf("⏳ Reading a mail from %s (%dsec)", from, counter*progressTimeout))
				} else {
					b.Edit(edit, fmt.Sprintf("⌛ Reading a mail from %s (%dsec)", from, counter*progressTimeout))
				}
			}
		}
	}()
}

func (r *Reader) startReadEmailBody(ctx context.Context, userId int64, seqNum int64) {
	imap := r.fact.ImapEmail()
	pnt := r.fact.PrintMsg()

	end := make(chan bool)
	defer func() {
		end <- true
		close(end)
	}()

	// Send start read
	r.startReadProgress(ctx, userId, seqNum, end)

	// Start read
	msg, err := imap.ReadEmail(ctx, r.imapUser, seqNum)
	if err != nil {
		logrus.Errorf("error read msg #%d: %v", seqNum, err)
		return
	}

	fmsg, err := pnt.PrintMsgWithBody(msg, r.imapUser.Login)
	if err != nil {
		logrus.Errorf("error print msg #%d: %v", seqNum, err)
		return
	}

	// Send result
	r.sendPrintMsg(fmsg, userId)
}

func (r *Reader) onButton(ctx context.Context, btnCtx bot.BtnContext) error {
	seqNum, err := strconv.ParseInt(btnCtx.Data(), 10, 32)
	if err != nil {
		logrus.Errorf("err parse string to int64: %v", err)
	}
	switch btnCtx.Unique() {
	case BtnMark:
	case BtnRead:
		go r.startReadEmailBody(ctx, btnCtx.UserId(), seqNum)
	default:
		logrus.Errorf("unknow btn type %s", btnCtx.Unique())
	}
	return nil
}

func (r *Reader) Start(ctx context.Context) {
	logrus.Infof("Start read unseen emails %s", r.imapUser.Login)
	imap := r.fact.ImapEmail()
	b := r.fact.Bot()
	pnt := r.fact.PrintMsg()
	storage := r.fact.Storage()

	emails, err := imap.ReadUnseenEmails(ctx, r.imapUser)
	if err != nil {
		logrus.Fatalln(err)
	}

	sort.Slice(emails, func(i, j int) bool {
		return emails[i].Date.Before(emails[j].Date)
	})

	for _, e := range emails {
		if err := storage.SaveMsgInfo(r.imapUser.Login, e); err != nil {
			logrus.Errorf("error save msg info: %v", err)
		}

		sid := strconv.FormatUint(uint64(e.SeqNum), 10)
		msg := pnt.PrintMsgEnvelope(e)
		for _, id := range r.userIds {
			if contains, err := storage.MsgWasSentToBotUser(r.imapUser.Login, e.SeqNum, id); err != nil {
				logrus.Errorf("error get msg contains from storage: %v", err)
			} else if contains {
				continue
			}

			inline := bot.NewInline(2, func(bc bot.BtnContext) error {
				return r.onButton(ctx, bc)
			})
			inline.Add("📩 Mark as read", BtnMark, sid)
			inline.Add("📧 Read", BtnRead, sid)
			if _, err := b.Send(id, msg, inline); err != nil {
				logrus.Errorf("error send msg: %v", err)
			} else {
				if err := storage.SaveMsgSentToBotUser(r.imapUser.Login, e.SeqNum, id); err != nil {
					logrus.Errorf("error save msg id to storage: %v", err)
				}
			}
		}
	}
}
