package email

import "context"

type ReadEmail interface {
	ReadUnseenEmails(ctx context.Context, user *ImapUser) ([]*MessageEnvelope, error)
	ReadEmail(ctx context.Context, user *ImapUser, msgUID uint32) (*Message, error)
}
