package types

type UID uint64

type EmailUser struct {
	ImapServer string
	Login      string
	Password   string
}
