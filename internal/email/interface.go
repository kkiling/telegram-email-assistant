package email

import "context"

type ReadEmail interface {
	ReadUnseenEmails(ctx context.Context, user *ImapUser) ([]*MessageEnvelope, error)
	ReadEmail(ctx context.Context, user *ImapUser, seqNum int64) (*Message, error)
}
