package email

import "context"

type ReadEmail interface {
	ReadUnseenEmails(ctx context.Context, user *ImapUser) ([]Message, error)
	ReadEmailBody(ctx context.Context, user *ImapUser, msgUID uint32) (*MessageWithBody, error)
}
