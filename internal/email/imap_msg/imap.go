package imapmsg

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/kiling91/telegram-email-assistant/internal/common"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/sirupsen/logrus"
)

type service struct {
	fact factory.Factory
}

func NewReadEmail(fact factory.Factory) email.ReadEmail {
	imap.CharsetReader = charset.Reader
	return &service{
		fact: fact,
	}
}

func (s *service) login(user *email.ImapUser) (*client.Client, error) {
	// Connect to server
	c, err := client.DialTLS(user.ImapServer, nil)
	if err != nil {
		return nil, fmt.Errorf("error connect to imap server: %w", err)
	}

	// Login
	if err := c.Login(user.Login, user.Password); err != nil {
		return nil, fmt.Errorf("error login in imap server: %w", err)
	}

	return c, nil
}

func (s *service) getUnseenEmails(client *client.Client) ([]uint32, error) {
	_, err := client.Select("INBOX", true)
	if err != nil {
		return nil, fmt.Errorf("error select inbox: %w", err)
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	seqNums, err := client.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("error search mail: %w", err)
	}

	return seqNums, nil
}

func (s *service) readEmailEnvelope(client *client.Client, seqNums ...uint32) ([]*email.MessageEnvelope, error) {
	result := make([]*email.MessageEnvelope, 0)
	if len(seqNums) == 0 {
		return result, nil
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqNums...)

	items := []imap.FetchItem{imap.FetchEnvelope}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqSet, items, messages)
	}()

	for msg := range messages {
		from := msg.Envelope.From[0]
		to := msg.Envelope.To[0]
		result = append(result, &email.MessageEnvelope{
			SeqNum:      int64(msg.SeqNum),
			Date:        msg.Envelope.Date,
			Subject:     msg.Envelope.Subject,
			FromAddress: from.MailboxName + from.HostName,
			FromName:    from.PersonalName,
			ToAddress:   to.MailboxName + to.HostName,
			ToName:      to.PersonalName,
		})
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("error read email envelope: %w", err)
	}

	return result, nil
}

func (s *service) saveFile(fileName string, body io.Reader, user string, seqNum int64) (string, error) {
	cfg := s.fact.Config()
	newPath, err := common.CreateFolderForEmail(cfg.App.FileDirectory, user, seqNum)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(newPath, fileName)
	b, _ := ioutil.ReadAll(body)
	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		return "", fmt.Errorf("error write file %s with error %w", filePath, err)
	}

	return filePath, nil
}

func (s *service) processReadBody(_ context.Context, mr *mail.Reader, user string, seqNum int64) (*email.MessageBody, error) {

	msgBody := email.MessageBody{
		TextHtml:        "",
		TextPlain:       "",
		InlineFiles:     make([]*email.InlineFile, 0),
		AttachmentFiles: make([]*email.AttachmentFile, 0),
	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error read email body: %w", err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, err := h.ContentType()
			if err != nil {
				return nil, err
			}

			switch contentType {
			case "text/plain":
				b, _ := ioutil.ReadAll(p.Body)
				msgBody.TextPlain = string(b)
			case "text/html":
				b, _ := ioutil.ReadAll(p.Body)
				msgBody.TextHtml = string(b)
			default:
				contentDisposition, contentDispositionParams, _ := h.ContentDisposition()
				if contentDisposition == "inline" {
					// This is an inline
					fileName := contentDispositionParams["filename"]
					attachmentId := h.Get("X-Attachment-Id")
					if attachmentId == "" {
						attachmentId = common.GetContentId(h.Get("Content-Id"))
						fileName = attachmentId
					}

					if attachmentId == "" {
						logrus.Errorf("seqNum: %d - inline attachmentId is empty", seqNum)
					} else {
						filePath, err := s.saveFile(attachmentId, p.Body, user, seqNum)
						if err != nil {
							return nil, err
						}
						msgBody.InlineFiles = append(msgBody.InlineFiles, &email.InlineFile{
							FileName:     fileName,
							FilePath:     filePath,
							AttachmentId: attachmentId,
						})
					}
				} else {
					logrus.Errorf("Unknown contentDisposition: %s", contentDisposition)
					logrus.Errorf("Unknown contentType: %s", contentType)
				}
			}
		case *mail.AttachmentHeader:
			// This is an attachment
			fileName, _ := h.Filename()

			if fileName == "" {
				fileName = common.GetContentId(h.Get("Content-Id"))
			}

			if fileName == "" {
				logrus.Warnf("seqNum: %d - attachment fileName is empty", seqNum)
			} else {
				filePath, err := s.saveFile(fileName, p.Body, user, seqNum)
				if err != nil {
					return nil, err
				}
				msgBody.AttachmentFiles = append(msgBody.AttachmentFiles, &email.AttachmentFile{
					FileName: fileName,
					FilePath: filePath,
				})
			}
		}
	}

	return &msgBody, nil
}

func (s *service) processReadEnvelope(seqNum int64, mr *mail.Reader) (*email.MessageEnvelope, error) {
	msgEnvelope := email.MessageEnvelope{
		SeqNum: seqNum,
	}

	// Print some info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		msgEnvelope.Date = date
	} else {
		return nil, fmt.Errorf("error get 'Date' from header: %w", err)
	}

	if from, err := header.AddressList("From"); err == nil {
		msgEnvelope.FromAddress = from[0].Address
		msgEnvelope.FromName = from[0].Name
	} else {
		return nil, fmt.Errorf("error get 'From' from header: %w", err)
	}

	if to, err := header.AddressList("To"); err == nil {
		msgEnvelope.ToAddress = to[0].Address
		msgEnvelope.ToName = to[0].Name
	} else {
		return nil, fmt.Errorf("error get 'Address' from header: %w", err)
	}

	if subject, err := header.Subject(); err == nil {
		msgEnvelope.Subject = subject
	} else {
		return nil, fmt.Errorf("error get 'Subject' from header: %w", err)
	}

	return &msgEnvelope, nil
}

func (s *service) readEmailBody(ctx context.Context, client *client.Client, user string, seqNum int64) (*email.Message, error) {
	// Select INBOX
	mbox, err := client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("error select mailbox: %w", err)
	}

	// Get the last message
	if mbox.Messages == 0 {
		return nil, fmt.Errorf("no message in mailbox")
	}

	// Select msg by uid
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uint32(seqNum))

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)

	if err := client.Fetch(seqSet, items, messages); err != nil {
		return nil, fmt.Errorf("error fetch email: %w", err)
	}

	msg := <-messages
	if msg == nil {
		return nil, fmt.Errorf("server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		return nil, fmt.Errorf("server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		return nil, fmt.Errorf("error create reader: %w", err)
	}

	msgEnvelope, err := s.processReadEnvelope(int64(msg.SeqNum), mr)
	if err != nil {
		return nil, err
	}

	msgBody, err := s.processReadBody(ctx, mr, user, seqNum)
	if err != nil {
		return nil, err
	}

	return &email.Message{
		Envelope: msgEnvelope,
		Body:     msgBody,
	}, nil
}

func (s *service) ReadUnseenEmails(_ context.Context, user *email.ImapUser) ([]*email.MessageEnvelope, error) {
	c, err := s.login(user)

	if err != nil {
		return nil, err
	}

	defer func(c *client.Client) {
		err := c.Logout()
		if err != nil {
			logrus.Errorf("error logout from imap server: %v", err)
		}
	}(c)

	// Select INBOX
	seqNums, err := s.getUnseenEmails(c)
	if err != nil {
		return nil, err
	}

	result, err := s.readEmailEnvelope(c, seqNums...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) ReadEmail(ctx context.Context, user *email.ImapUser, seqNum int64) (*email.Message, error) {
	c, err := s.login(user)
	defer func(c *client.Client) {
		err := c.Logout()
		if err != nil {
			logrus.Errorf("error logout from imap server: %v", err)
		}
	}(c)

	if err != nil {
		return nil, err
	}

	return s.readEmailBody(ctx, c, user.Login, seqNum)
}
