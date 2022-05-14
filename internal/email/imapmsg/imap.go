package imapmsg

import (
	"context"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/kiling91/telegram-email-assistant/internal/common"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
)

type service struct {
	fact factory.Factory
}

func NewReadEmail(fact factory.Factory) email.ReadEmail {
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
	UIDs, err := client.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("error search mail: %w", err)
	}

	return UIDs, nil
}

func (s *service) readEmailEnvelope(client *client.Client, UIDs ...uint32) ([]email.Message, error) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(UIDs...)

	items := []imap.FetchItem{imap.FetchEnvelope}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqSet, items, messages)
	}()

	result := make([]email.Message, 0)
	for msg := range messages {
		result = append(result, email.Message{
			Uid:     msg.SeqNum,
			Date:    msg.Envelope.Date,
			Subject: msg.Envelope.Subject,
			From:    msg.Envelope.From[0].MailboxName + msg.Envelope.From[0].HostName,
			To:      msg.Envelope.To[0].MailboxName + msg.Envelope.To[0].HostName,
		})
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("error read email envelope: %w", err)
	}

	return result, nil
}

func (s *service) saveFile(fileName string, body io.Reader, emailUser string, msgUID uint32) (string, error) {
	cfg := s.fact.Config()
	newPath, err := common.CreateFolderForEmailUser(cfg.FileStorageDir, emailUser, msgUID)
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

func (s *service) readEmailBody(client *client.Client, emailUser string, msgUID uint32) (*email.MessageWithBody, error) {
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
	seqSet.AddNum(msgUID)

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

	result := email.MessageWithBody{
		InlineFiles:     make([]*email.InlineFile, 0),
		AttachmentFiles: make([]*email.AttachmentFile, 0),
	}
	result.Uid = msg.SeqNum
	// Print some info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		result.Date = date
	} else {
		return nil, fmt.Errorf("error get 'Date' from header: %w", err)
	}

	if from, err := header.AddressList("From"); err == nil {
		result.From = from[0].Address
	} else {
		return nil, fmt.Errorf("error get 'From' from header: %w", err)
	}

	if to, err := header.AddressList("To"); err == nil {
		result.To = to[0].Address
	} else {
		return nil, fmt.Errorf("error get 'Address' from header: %w", err)
	}

	if subject, err := header.Subject(); err == nil {
		result.Subject = subject
	} else {
		return nil, fmt.Errorf("error get 'Subject' from header: %w", err)
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
				result.TextPlain = string(b)
			case "text/html":
				b, _ := ioutil.ReadAll(p.Body)
				result.TextHtml = string(b)
			default:
				contentDisposition, contentDispositionParams, _ := h.ContentDisposition()
				if contentDisposition == "inline" {
					// This is an inline
					fileName := contentDispositionParams["filename"]
					filePath, err := s.saveFile(fileName, p.Body, emailUser, msgUID)
					if err != nil {
						return nil, err
					}
					result.InlineFiles = append(result.InlineFiles, &email.InlineFile{
						FileName: fileName,
						FilePath: filePath,
					})
				} else {
					log.Printf("Unknown contentDisposition: %s", contentDisposition)
					log.Printf("Unknown contentType: %s", contentType)
				}
			}
		case *mail.AttachmentHeader:
			// This is an attachment
			fileName, _ := h.Filename()
			filePath, err := s.saveFile(fileName, p.Body, emailUser, msgUID)
			if err != nil {
				return nil, err
			}
			result.AttachmentFiles = append(result.AttachmentFiles, &email.AttachmentFile{
				FileName: fileName,
				FilePath: filePath,
			})
		}
	}

	return &result, nil
}

func (s *service) ReadUnseenEmails(ctx context.Context, user *email.ImapUser) ([]email.Message, error) {
	c, err := s.login(user)
	defer func(c *client.Client) {
		err := c.Logout()
		if err != nil {
			log.Println("error logout from imap server: %w", err)
		}
	}(c)

	if err != nil {
		return nil, err
	}

	// Select INBOX
	UIDs, err := s.getUnseenEmails(c)
	if err != nil {
		return nil, err
	}

	result, err := s.readEmailEnvelope(c, UIDs...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) ReadEmailBody(ctx context.Context, user *email.ImapUser, msgUID uint32) (*email.MessageWithBody, error) {
	c, err := s.login(user)
	defer func(c *client.Client) {
		err := c.Logout()
		if err != nil {
			log.Println("error logout from imap server: %w", err)
		}
	}(c)

	if err != nil {
		return nil, err
	}

	return s.readEmailBody(c, user.Login, msgUID)
}
