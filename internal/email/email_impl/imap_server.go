package email_impl

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/kiling91/telegram-email-assistant/internal/types"
	"io"
	"io/ioutil"
	"log"
)

type service struct {
	fact factory.Factory
}

func NewImapServer(fact factory.Factory) email.ImapServer {
	return &service{
		fact: fact,
	}
}

func (s *service) readEmailBody(client *client.Client, msgUid uint32) error {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(msgUid)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := client.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	msg := <-messages
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	// Print some info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		log.Println("Date:", date)
	}
	if from, err := header.AddressList("From"); err == nil {
		log.Println("From:", from)
	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			log.Println("Got text: %v", string(b))
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Println("Got attachment: %v", filename)
		}
	}

	return nil
}

func (s *service) login(user *types.EmailUser) (*client.Client, error) {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(user.ImapServer, nil)
	if err != nil {
		return nil, fmt.Errorf("error connect to imap server: %w", err)
	}
	log.Println("Connected")

	// Login
	if err := c.Login(user.Login, user.Password); err != nil {
		return nil, fmt.Errorf("error login in imap server: %w", err)
	}
	log.Println("Logged in")

	return c, nil
}

func (s *service) getUnseenEmails(client *client.Client) ([]uint32, error) {
	mbox, err := client.Select("INBOX", true)
	if err != nil {
		return nil, fmt.Errorf("error select inbox: %w", err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	UIDs, err := client.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("error search mail: %w", err)
	}

	return UIDs, nil
}

func (s *service) readEmailEnvelope(client *client.Client, UIDs ...uint32) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(UIDs...)

	items := []imap.FetchItem{imap.FetchEnvelope}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqSet, items, messages)
	}()

	log.Println("Unseen messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}

func (s *service) ReadUnseenEmails(user *types.EmailUser) error {
	c, err := s.login(user)
	defer func(c *client.Client) {
		err := c.Logout()
		if err != nil {
			log.Println("error logout from imap server: %w", err)
		}
	}(c)

	if err != nil {
		return err
	}

	// Select INBOX
	UIDs, err := s.getUnseenEmails(c)
	if err != nil {
		log.Fatal(err)
	}

	s.readEmailEnvelope(c, UIDs...)

	return nil
}
